package server

import (
	"github.com/google/uuid"
	"github.com/mattwalo32/RealTimeAPI/internal/conn"
	"github.com/mattwalo32/RealTimeAPI/internal/messages"
	"github.com/mattwalo32/RealTimeAPI/internal/timer"
	log "github.com/sirupsen/logrus"
	"net"
	"sync"
)

const (
	UDP_RECEIVING_CHAN_SIZE = 15
	MIN_RECEIVING_CHAN_CAP  = 2
)

type MessageRouter struct {
	// Maps client ID to client data
	clients map[uuid.UUID]*Client

	// Maps room ID to room struct
	rooms map[uuid.UUID]*Room

	// Maps room type to configuration
	roomConfigs map[int]*RoomConfig

	// Maps messageID to outstanding message retry event IDs
	messageRetryEventIDs map[uuid.UUID]uuid.UUID

	// The current packet count. Used to sequentially number packets.
	packetCount int

	// Timer object used for resending messages
	timer *timer.Timer

	// Config passed in constructor
	config *MessageRouterConfig

	// Used to send packets over UDP
	udpManager *conn.UDPManager

	// Packets are written to this channel by the UDPManager
	udpReceivingChan chan conn.Packet

	doneChan chan bool
	lock     sync.Mutex
}

type Client struct {
	Address           net.UDPAddr
	ID                uuid.UUID
	RoomID            uuid.UUID
	AppData           string
	lastContactTimeMs int64
}

type MessageRouterConfig struct {
	// Passed to UDPManager
	Address string

	// Max amount of times to retry sending a reliable message
	MaxMessageRetries int

	// Time between reliable message retries
	MessageRetryTimeoutMs uint64

	// How often to heartbeat with client to check if alive
	HeartbeatIntervalMs uint64

	// A client will be disconnected if they don't respond in this many heartbeat's time
	HeartbeatActivationMultiplier float64
}

func NewMessageRouter(config MessageRouterConfig) *MessageRouter {
	udpReceivingChan := make(chan conn.Packet, UDP_RECEIVING_CHAN_SIZE)

	udpConfig := conn.UDPManagerConfig{
		ReceivingChan: udpReceivingChan,
		Address:       config.Address,
	}

	router := &MessageRouter{
		udpManager:           conn.NewUDPManager(udpConfig),
		udpReceivingChan:     udpReceivingChan,
		doneChan:             make(chan bool),
		clients:              make(map[uuid.UUID]*Client),
		rooms:                make(map[uuid.UUID]*Room),
		roomConfigs:          make(map[int]*RoomConfig),
		messageRetryEventIDs: make(map[uuid.UUID]uuid.UUID),
		timer:                timer.NewTimer(),
		config:               &config,
		packetCount:          0,
	}

	// Remove this later. Only for debugging
	log.SetLevel(log.DebugLevel)

	go router.decodeMessages()
	return router
}

/**
  * Registers a room config wth the message router. This should be called as early as possible in the server's startup.
  * When a "FIND_ROOM" message searching for a room of this type is received, a new room will be created off this config if needed.
  * @param config: The room configuration from which all rooms of this type are made
  * @param roomType: A user-defined ID for this type of room. 
  */
func (router *MessageRouter) RegisterRoomConfiguration(roomConfig RoomConfig, roomType int) {
	_,isIDTaken := router.roomConfigs[roomType]
	if isIDTaken {
		log.WithFields(log.Fields{
			"RoomType": roomType,
		}).Error("Attempted to register room config with an existing room type. All room configs must have a unique room type! Room config will not be registered.")
		return
	}

	router.roomConfigs[roomType] = &roomConfig
}

func (router *MessageRouter) decodeMessages() {
	var udpMsg conn.Packet

	for {
		select {
		case <-router.doneChan:
			return
		case udpMsg = <-router.udpReceivingChan:
		}

		message, err := messages.DecodeFromHeader(udpMsg.Data)
		if err != nil {
			log.Warn(err)
			continue
		}

		log.WithFields(log.Fields{
			"Type": message.GetMessageType(),
		}).Debug("Decoded Message")
		message.SetSource(udpMsg.Address)
		router.processMessage(message)
	}
}

func (router *MessageRouter) SendMessageUnreliably(msg messages.Message) {
	router.lock.Lock()
	defer router.lock.Unlock()
	msg.SetResponseRequired(false)
	msg.SetPacketNumber(router.packetCount)
	router.sendMessage(msg)
}

func (router *MessageRouter) sendMessage(msg messages.Message) {
	data, err := messages.EncodeWithHeader(msg)
	if err != nil {
		log.Warn(err)
		return
	}

	udpMsg := conn.Packet{
		Data:    data,
		Address: msg.GetDestination(),
	}

	log.WithFields(log.Fields{
		"Type":        msg.GetMessageType(),
		"Destination": msg.GetDestination(),
	}).Debug("Sending Message")
	router.udpManager.SendPacket(udpMsg)
	router.packetCount++
}

func (router *MessageRouter) SendMessageReliably(msg messages.Message) {
	router.lock.Lock()
	defer router.lock.Unlock()
	msg.SetResponseRequired(true)
	msg.SetPacketNumber(router.packetCount)
	router.createTimerForMessage(msg)
	router.sendMessage(msg)
}

func (router *MessageRouter) createTimerForMessage(msg messages.Message) {
	c := router.config
	id := router.timer.AddRepeatingEvent(router.onMessageRetry, msg, c.MessageRetryTimeoutMs, c.MaxMessageRetries)
	router.messageRetryEventIDs[msg.GetID()] = id
	log.WithFields(log.Fields{
		"ID": msg.GetID(),
	}).Debug("Adding timer for message")
}

func (router *MessageRouter) onMessageRetry(message interface{}) {
	log.Debug("Resending message")
	msg := message.(messages.Message)
	router.sendMessage(msg)
}

func (router *MessageRouter) Stop() {
	router.udpManager.Stop()
	router.timer.Stop()
	close(router.doneChan)
}

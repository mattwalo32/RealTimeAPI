package server

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/mattwalo32/RealTimeAPI/internal/conn"
	"github.com/mattwalo32/RealTimeAPI/internal/messages"
	"github.com/mattwalo32/RealTimeAPI/internal/timer"
	log "github.com/sirupsen/logrus"
	"sync"
	"net"
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

	// Maps messageID to outstanding message retry event IDs
	messageRetryEventIDs map[uuid.UUID]uuid.UUID

	// The current packet count. Used to sequentially number packets.
	packetCount          int

	// Timer object used for resending messages
	timer                *timer.Timer

	// Config passed in constructor
	config               *MessageRouterConfig

	// Used to send packets over UDP
	udpManager           *conn.UDPManager

	// Packets are written to this channel by the UDPManager
	udpReceivingChan     chan conn.Packet

	doneChan             chan bool
	lock sync.Mutex
}

type Client struct {
	Address     net.UDPAddr
	ID          uuid.UUID
	AppData     string
	lastContactTimeMs uint64
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

	handler := &MessageRouter{
		udpManager:       conn.NewUDPManager(udpConfig),
		udpReceivingChan: udpReceivingChan,
		doneChan:         make(chan bool),
		clients: make(map[uuid.UUID]*Client),
		rooms: make(map[uuid.UUID]*Room),
		messageRetryEventIDs: make(map[uuid.UUID]uuid.UUID),
		timer:            timer.NewTimer(),
		config:           &config,
		packetCount:      0,
	}

	// Remove this later. Only for debugging
	log.SetLevel(log.DebugLevel)

	go handler.decodeMessages()
	return handler
}

func (handler *MessageRouter) decodeMessages() {
	var udpMsg conn.Packet

	for {
		select {
		case <-handler.doneChan:
			return
		case udpMsg = <-handler.udpReceivingChan:
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
		handler.processMessage(message)
	}
}

func (handler *MessageRouter) SendMessageUnreliably(msg messages.Message) {
	handler.lock.Lock()
	defer handler.lock.Unlock()
	msg.SetResponseRequired(false)
	msg.SetPacketNumber(handler.packetCount)
	handler.sendMessage(msg)
}

func (handler *MessageRouter) sendMessage(msg messages.Message) {
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
		"Type": msg.GetMessageType(),
		"Destination": msg.GetDestination(),
	}).Debug("Sending Message")
	handler.udpManager.SendPacket(udpMsg)
	handler.packetCount++
}

func (handler *MessageRouter) SendMessageReliably(msg messages.Message) {
	handler.lock.Lock()
	defer handler.lock.Unlock()
	msg.SetResponseRequired(true)
	msg.SetPacketNumber(handler.packetCount)
	handler.createTimerForMessage(msg)
	handler.sendMessage(msg)
}

func (handler *MessageRouter) createTimerForMessage(msg messages.Message) {
	c := handler.config
	id := handler.timer.AddRepeatingEvent(handler.onMessageRetry, msg, c.MessageRetryTimeoutMs, c.MaxMessageRetries)
	handler.messageRetryEventIDs[msg.GetID()] = id
	log.WithFields(log.Fields{
		"ID": msg.GetID(),
	}).Debug("Adding timer for message")
}

func (handler *MessageRouter) onMessageRetry(message interface{}) {
	log.Debug("Resending message")
	msg := message.(messages.Message)
	handler.sendMessage(msg)
}

func (handler *MessageRouter) Stop() {
	handler.udpManager.Stop()
	handler.timer.Stop()
	close(handler.doneChan)
}

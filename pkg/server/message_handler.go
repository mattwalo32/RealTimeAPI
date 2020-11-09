package server

import (
	"fmt"
	"net"
	"github.com/google/uuid"
	"github.com/mattwalo32/RealTimeAPI/internal/conn"
	"github.com/mattwalo32/RealTimeAPI/internal/messages"
	log "github.com/sirupsen/logrus"
)

const (
	UDP_RECEIVING_CHAN_SIZE = 5
	MIN_RECEIVING_CHAN_CAP  = 2
)

type MessageHandler struct {
	// Maps client ID to their data
	clients          map[uuid.UUID]*ClientData

	// Maps messageID to outstanding message retry event IDs
	messageRetryEventIDs map[uuid.UUID]uuid.UUID
	packetCount      int
	config           *MessageHandlerConfig
	udpManager       *conn.UDPManager
	udpReceivingChan chan conn.Message
	doneChan         chan bool
}

type ClientData struct {
	lastContactTimeMs uint64
	immutableData *ImmutableClientData
}

type ImmutableClientData struct {
	Address net.UDPAddr
	ID uuid.UUID
	ProfileData string
}

type MessageHandlerConfig struct {
	// Acts as callback mechanism for decoded messages
	MessageReceivingChan chan messages.Encodable

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

func NewMessageHandler(config MessageHandlerConfig) *MessageHandler {
	udpReceivingChan := make(chan conn.Message, UDP_RECEIVING_CHAN_SIZE)
	err := assertConfigValid(&config)
	if err != nil {
		log.Fatal(err)
	}

	udpConfig := conn.UDPManagerConfig{
		ReceivingChan: udpReceivingChan,
		Address:       config.Address,
	}

	handler := &MessageHandler{
		udpManager:       conn.NewUDPManager(udpConfig),
		udpReceivingChan: udpReceivingChan,
		doneChan:         make(chan bool),
		config:           &config,
		packetCount: 0,
	}

	go handler.decodeMessages()
	return handler
}

func assertConfigValid(config *MessageHandlerConfig) error {
	cap := cap(config.MessageReceivingChan)
	if cap < MIN_RECEIVING_CHAN_CAP {
		return fmt.Errorf("Recieving message channel must have capacity %v, got: %v", MIN_RECEIVING_CHAN_CAP, cap)
	}

	return nil
}

func (handler *MessageHandler) decodeMessages() {
	var udpMsg conn.Message

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

		message.SetSource(udpMsg.Address)
		handler.processMessage(message)
	}
}

func (handler *MessageHandler) SendMessageUnreliably(msg messages.Encodable) {
	msg.SetResponseRequired(false)
	msg.SetPacketNumber(handler.packetCount)
	handler.sendMessage(msg)
}

func (handler *MessageHandewr) sendMessage(msg messages.Encodable) {
	data, err := messages.EncodeWithHeader(msg)
	if err != nil {
		log.Warn(err)
		return
	}

	udpMsg := conn.Message{
		Data:    data,
		Address: msg.GetDestination(),
	}

	handler.udpManager.SendMessage(udpMsg)
	handler.packetCount++
}
func (handler *MessageHandler) SendMessageReliably(msg messages.Encodable) {
	msg.SetResponseRequired(true)
	msg.SetPacketNumber(handler.packetCount)
	handler.setTimerForMessage(msg.GetID())
	handler.sendMessage(msg)
}

func (hander *MessageHandler) setTimerForMessage(msgID uuid.UUID) {
	// TODO:
}

func (handler *MessageHandler) Stop() {
	handler.udpManager.Stop()
	close(handler.doneChan)
}

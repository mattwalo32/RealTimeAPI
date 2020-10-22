package server

import (
	"fmt"
	"github.com/mattwalo32/RealTimeAPI/internal/conn"
	"github.com/mattwalo32/RealTimeAPI/internal/messages"
	log "github.com/sirupsen/logrus"
)

const (
	UDP_RECEIVING_CHAN_SIZE = 5
	MIN_RECEIVING_CHAN_CAP  = 2
)

type MessageHandler struct {
	config           *MessageHandlerConfig
	udpManager       *conn.UDPManager
	udpReceivingChan chan conn.Message
	doneChan         chan bool
}

type MessageHandlerConfig struct {
	// Acts as callback mechanism for decoded messages
	MessageReceivingChan chan messages.Encodable

	// Passed to UDPManager
	Address string
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
		handler.config.MessageReceivingChan <- message

		// TODO: Don't send back acknowledgment messages
	}
}

func (handler *MessageHandler) SendMessage(msg messages.Encodable) {
	data, err := msg.Encode()
	if err != nil {
		log.Warn(err)
		return
	}

	udpMsg := conn.Message{
		Data:    data,
		Address: msg.GetDestination(),
	}

	handler.udpManager.SendMessage(udpMsg)
}

func (handler *MessageHandler) SendMessageReliably(msg *messages.Encodable) {
	// TODO: Handle responses and timing
}

func (handler *MessageHandler) Stop() {
	handler.udpManager.Stop()
	close(handler.doneChan)
}

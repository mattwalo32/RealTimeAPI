package server

import (
	"github.com/mattwalo32/RealTimeAPI/internal/messages"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func (handler *MessageHandler) processMessage(msg messages.Encodable) {
	switch msg.GetMessageType() {
	case messages.MESSAGE_ACKNOWLEDGE:
		handler.processAcknowledge(msg.(*messages.AcknowledgementMessage))
	case messages.MESSAGE_JOIN_SERVER:
		// TODO: Send back uuid
	default:
		handler.config.MessageReceivingChan <- msg
	}

	if msg.IsResponseRequired() {
		handler.acknowledgeMessage(msg)
	}
}

func (handler *MessageHandler) processAcknowledge(msg *messages.AcknowledgementMessage) {
	handler.removeMessageTimer(msg)
}

func (handler *MessageHandler) removeMessageTimer(msg *messages.AcknowledgementMessage) {
	handler.lock.Lock()
	defer handler.lock.Unlock()
	evtID, evtExists := handler.messageRetryEventIDs[msg.AcknowledgedMessageID]
	if !evtExists {
		log.WithFields(log.Fields{
			"ID": msg.AcknowledgedMessageID,
		}).Debug("Acknowledged message without a timer")
		return
	}

	log.WithFields(log.Fields{
		"ID": msg.AcknowledgedMessageID,
	}).Debug("Recieved acknowledgement for message")
	handler.timer.RemoveEvent(evtID)
	delete(handler.messageRetryEventIDs, msg.AcknowledgedMessageID)
}

func (handler *MessageHandler) acknowledgeMessage(msg messages.Encodable) {
	ackMessage := &messages.AcknowledgementMessage{
		SourceAddr: *handler.udpManager.GetUDPAddr(),
		DestAddr: msg.GetSource(),
		MessageID: uuid.New(),
		AcknowledgedMessageID: msg.GetID(),
	}

	handler.SendMessageUnreliably(ackMessage)
}
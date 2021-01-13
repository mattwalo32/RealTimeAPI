package server

import (
	"github.com/mattwalo32/RealTimeAPI/internal/messages"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func (handler *MessageRouter) processMessage(msg messages.Message) {
	switch msg.GetMessageType() {
	case messages.MESSAGE_ACKNOWLEDGE:
		handler.processAcknowledge(msg.(*messages.AcknowledgementMessage))
	case messages.MESSAGE_JOIN_SERVER:
		handler.processJoinServer(msg.(*messages.JoinServerMessage))
	default:
	}

	// TODO: Don't allow acknowledgement messages to get acknowledged
	if msg.IsResponseRequired() {
		handler.acknowledgeMessage(msg)
	}
}

func (handler *MessageRouter) removeMessageTimer(msg *messages.AcknowledgementMessage) {
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

func (handler *MessageRouter) acknowledgeMessage(msg messages.Message) {
	ackMessage := &messages.AcknowledgementMessage{
		SourceAddr: *handler.udpManager.GetUDPAddr(),
		DestAddr: msg.GetSource(),
		MessageID: uuid.New(),
		AcknowledgedMessageID: msg.GetID(),
	}

	handler.SendMessageUnreliably(ackMessage)
}
package server

import (
	"github.com/mattwalo32/RealTimeAPI/internal/messages"
)

func (handler *MessageHandler) processMessage(msg messages.Encodable) {
	switch msg.GetMessageType() {
	case messages.MESSAGE_ACKNOWLEDGE:
		// TODO: Handle acknowledgment
	case messages.MESSAGE_JOIN_SERVER:
		// TODO: Send back uuid
	default:
		handler.config.MessageReceivingChan <- msg
	}
}

func (handler *MessageHandler) processAcknowledge(msg messages.AcknowledgementMessage) {
	handler.removeMessageTimer(msg)
}

func (handler *MessageHandler) removeMessageTimer(msg messages.AcknowledgementMessage) {
	evtID, evtExists := handler.messageRetryEventIDs[msg.GetID()]
	if !evtExists {
		return
	}

	handler.timer.RemoveEvent(evtID)
	delete(handler.messageRetryEventIDs, evtID)
}
package server

import (
	"github.com/mattwalo32/RealTimeAPI/internal/messages"
)

func (handler *MessageHandler) processAcknowledge(msg *messages.AcknowledgementMessage) {
	handler.removeMessageTimer(msg)
}
package server

import (
	"github.com/mattwalo32/RealTimeAPI/internal/messages"
)

func (handler *MessageRouter) processAcknowledge(msg *messages.AcknowledgementMessage) {
	handler.removeMessageTimer(msg)
}
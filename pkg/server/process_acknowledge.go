package server

import (
	"github.com/mattwalo32/RealTimeAPI/internal/messages"
)

func (router *MessageRouter) processAcknowledge(msg *messages.AcknowledgementMessage) {
	router.removeMessageTimer(msg)
}

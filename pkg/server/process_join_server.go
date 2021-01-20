package server

import (
	"github.com/mattwalo32/RealTimeAPI/internal/messages"
	"github.com/mattwalo32/RealTimeAPI/internal/util"
	"github.com/google/uuid"
)

func (router *MessageRouter) processJoinServer(msg *messages.JoinServerMessage) {
	client := &Client{
		Address: msg.GetSource(),
		ID: uuid.New(),
		RoomID: uuid.Nil,
		AppData: msg.AppData,
		lastContactTimeMs: util.CurrentTimestampMs(),
	}

	// TODO: Track client heartbeats
	router.clients[client.ID] = client

	// TODO: Send status back
}

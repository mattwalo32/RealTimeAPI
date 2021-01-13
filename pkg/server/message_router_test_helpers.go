package server

import (
	"github.com/google/uuid"
)

/**
 * Creates a fake client with the given ID and assigns it the given room ID.
 * For testing only
 **/
func (router *MessageRouter) createMockClientInRoom(clientID uuid.UUID, roomID uuid.UUID) {
	mockClient := &Client{
		ID:     clientID,
		RoomID: roomID,
	}

	router.clients[clientID] = mockClient
}

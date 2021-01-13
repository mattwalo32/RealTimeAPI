package server;

import (
	"github.com/mattwalo32/RealTimeAPI/internal/messages"
	"github.com/google/uuid"
)

type RoomEventHandler interface {
	// Called whenever a client joins a room. It is up to the room to track the number of
	// players in it and set isRoomOpen to 'false' when full.
	OnClientConnected(Client)

	// Called when a client voluntarily leaves or loses connection
	OnClientDisconnected(Client)

	// Called when a message for any of the room's cilents are received
	OnMessageRecieved(messages.Message)
}

type Room struct {
	// Passed in constructor
	config RoomConfig

	// If a room is open, players will be put into it
	isRoomOpen bool

	// List of all clients in room 
	clients []*Client

	ID uuid.UUID
}

type RoomConfig struct {
	// Defines how the room should react to events
	EventHandler RoomEventHandler

	// If true, non-game messages (such as acknowledge messages) will be routed to the room. Recommended to be false.
	ShouldReceiveVerboseMessages bool

	// All rooms for the same app should have the same app ID
	ApplicationID uuid.UUID

	// The expected capacity of the room. This is used purely for performance, the capcity can grow or be unutilized
	InitialCapacity int
}

func NewRoom(config RoomConfig) *Room {
	room := &Room{
		config: config,
		isRoomOpen: true,
		clients: make([]*Client, 0, config.InitialCapacity),
		ID: uuid.New(),
	}

	return room
}
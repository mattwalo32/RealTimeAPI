package server

import (
	"github.com/mattwalo32/RealTimeAPI/internal/messages"
	"github.com/google/uuid"
)

type messageQueue struct {
	messages []messages.Message
}

func createMockRoom() *Room {
	evtHandler := &messageQueue{}

	config := RoomConfig{
		EventHandler: evtHandler,
		ShouldReceiveVerboseMessages: true,
		ApplicationID: uuid.New(),
		InitialCapacity: 10,
	}

	return NewRoom(config)
}

func (queue *messageQueue) OnClientConnected(client Client) {}

func (queue *messageQueue) OnClientDisconnected(client Client) {}

func (queue *messageQueue) OnMessageRecieved(msg messages.Message) {
	queue.messages = append(queue.messages, msg)
}

func (queue *messageQueue) DequeueMessage() messages.Message {
	queueLen := len(queue.messages)
	if queueLen < 1 {
		return nil
	}

	msg := queue.messages[0]

	if queueLen > 1 {
		copy(queue.messages[1:], queue.messages[0:queueLen-1])
	}

	return msg
}

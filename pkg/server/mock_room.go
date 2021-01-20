package server

import (
	"github.com/google/uuid"
	"github.com/mattwalo32/RealTimeAPI/internal/messages"
)

type messageQueueHandler struct {
	messages []messages.Message
}

func newMockRoomConfig() RoomConfig {
	evtHandler := &messageQueueHandler{}

	config := RoomConfig{
		EventHandler:                 evtHandler,
		ShouldReceiveVerboseMessages: true,
		ApplicationID:                uuid.New(),
		InitialCapacity:              10,
	}

	return config
}

func (queue *messageQueueHandler) OnClientConnected(client Client) {}

func (queue *messageQueueHandler) OnClientDisconnected(client Client) {}

func (queue *messageQueueHandler) OnMessageRecieved(msg messages.Message) {
	queue.messages = append(queue.messages, msg)
}

func (queue *messageQueueHandler) DequeueMessage() messages.Message {
	queueLen := len(queue.messages)
	if queueLen < 1 {
		return nil
	}

	msg := queue.messages[0]

	if queueLen > 1 {
		queue.messages = queue.messages[1:]
	} else {
		queue.messages = []messages.Message{}
	}

	return msg
}

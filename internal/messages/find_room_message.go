package messages

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"encoding/json"
)

type FindRoomMessage struct {
	UserID uuid.UUID
	ShouldStartWhenFull bool
	MinPlayers int
	MaxPlayers int
}

func (msg *FindRoomMessage) GetMessageType() int {
	return MESSAGE_FIND_ROOM
}

func (msg *FindRoomMessage) Encode() ([]byte, error) {
	return json.Marshal(msg)
}

func (msg *FindRoomMessage) Decode(data []byte) {
	err := json.Unmarshal(data, msg)
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Warn("Error decoding FindRoomMessage")
	}
}
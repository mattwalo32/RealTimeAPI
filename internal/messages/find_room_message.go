package messages

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"encoding/json"
	"net"
)

type FindRoomMessage struct {
	SourceAddr *net.UDPAddr
	DestAddr *net.UDPAddr

	UserID uuid.UUID
	ShouldStartWhenFull bool
	MinPlayers int
	MaxPlayers int
}

func (msg *FindRoomMessage) GetSource() *net.UDPAddr {
	return msg.SourceAddr
}

func (msg *FindRoomMessage) SetSource(addr *net.UDPAddr) {
	msg.SourceAddr = addr
}

func (msg *FindRoomMessage) SetDestination(addr *net.UDPAddr) {
	msg.DestAddr = addr
}

func (msg *FindRoomMessage) GetDestination() *net.UDPAddr {
	return msg.DestAddr
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
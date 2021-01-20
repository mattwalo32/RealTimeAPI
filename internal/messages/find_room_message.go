package messages

import (
	"encoding/json"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net"
)

type FindRoomMessage struct {
	SourceAddr       net.UDPAddr
	DestAddr         net.UDPAddr
	MessageID        uuid.UUID
	PacketNumber     int
	ResponseRequired bool

	UserID              uuid.UUID
	ClientID            uuid.UUID
	ShouldStartWhenFull bool
	RoomType            int
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

func (msg *FindRoomMessage) SetResponseRequired(isRequired bool) {
	msg.ResponseRequired = isRequired
}

func (msg *FindRoomMessage) IsResponseRequired() bool {
	return msg.ResponseRequired
}

func (msg *FindRoomMessage) SetPacketNumber(count int) {
	msg.PacketNumber = count
}

func (msg *FindRoomMessage) GetPacketNumber() int {
	return msg.PacketNumber
}

func (msg *FindRoomMessage) GetID() uuid.UUID {
	return msg.MessageID
}

func (msg *FindRoomMessage) GetSource() net.UDPAddr {
	return msg.SourceAddr
}

func (msg *FindRoomMessage) SetSource(addr net.UDPAddr) {
	msg.SourceAddr = addr
}

func (msg *FindRoomMessage) SetDestination(addr net.UDPAddr) {
	msg.DestAddr = addr
}

func (msg *FindRoomMessage) GetDestination() net.UDPAddr {
	return msg.DestAddr
}

func (msg *FindRoomMessage) GetMessageType() int {
	return MESSAGE_FIND_ROOM
}

func (msg *FindRoomMessage) GetClientID() uuid.UUID {
	return msg.ClientID
}

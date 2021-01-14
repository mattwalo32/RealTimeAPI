package messages

import (
	"encoding/json"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net"
)

type JoinServerMessage struct {
	SourceAddr       net.UDPAddr
	DestAddr         net.UDPAddr
	MessageID        uuid.UUID
	PacketNumber     int
	ResponseRequired bool

	ApplicationID    uuid.UUID
	AppData          string
}

func (msg *JoinServerMessage) Encode() ([]byte, error) {
	return json.Marshal(msg)
}

func (msg *JoinServerMessage) Decode(data []byte) {
	err := json.Unmarshal(data, msg)
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Warn("Error decoding JoinServerMessage")
	}
}

func (msg *JoinServerMessage) SetResponseRequired(isRequired bool) {
	msg.ResponseRequired = isRequired
}

func (msg *JoinServerMessage) IsResponseRequired() bool {
	return msg.ResponseRequired
}

func (msg *JoinServerMessage) SetPacketNumber(count int) {
	msg.PacketNumber = count
}

func (msg *JoinServerMessage) GetPacketNumber() int {
	return msg.PacketNumber
}

func (msg *JoinServerMessage) GetID() uuid.UUID {
	return msg.MessageID
}

func (msg *JoinServerMessage) GetSource() net.UDPAddr {
	return msg.SourceAddr
}

func (msg *JoinServerMessage) SetSource(addr net.UDPAddr) {
	msg.SourceAddr = addr
}

func (msg *JoinServerMessage) SetDestination(addr net.UDPAddr) {
	msg.DestAddr = addr
}

func (msg *JoinServerMessage) GetDestination() net.UDPAddr {
	return msg.DestAddr
}

func (msg *JoinServerMessage) GetMessageType() int {
	return MESSAGE_JOIN_SERVER
}

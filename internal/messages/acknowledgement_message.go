package messages

import (
	"encoding/json"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net"
)

type AcknowledgementMessage struct {
	SourceAddr       net.UDPAddr
	DestAddr         net.UDPAddr
	MessageID        uuid.UUID
	PacketNumber     int
	ResponseRequired bool

	AcknowledgedMessageID uuid.UUID
}

func (msg *AcknowledgementMessage) Encode() ([]byte, error) {
	return json.Marshal(msg)
}

func (msg *AcknowledgementMessage) Decode(data []byte) {
	err := json.Unmarshal(data, msg)
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Warn("Error decoding AcknowledgementMessage")
	}
}

func (msg *AcknowledgementMessage) SetResponseRequired(isRequired bool) {
	msg.ResponseRequired = isRequired
}

func (msg *AcknowledgementMessage) IsResponseRequired() bool {
	return msg.ResponseRequired
}

func (msg *AcknowledgementMessage) SetPacketNumber(count int) {
	msg.PacketNumber = count
}

func (msg *AcknowledgementMessage) GetPacketNumber() int {
	return msg.PacketNumber
}

func (msg *AcknowledgementMessage) GetID() uuid.UUID {
	return msg.MessageID
}

func (msg *AcknowledgementMessage) GetSource() net.UDPAddr {
	return msg.SourceAddr
}

func (msg *AcknowledgementMessage) SetSource(addr net.UDPAddr) {
	msg.SourceAddr = addr
}

func (msg *AcknowledgementMessage) SetDestination(addr net.UDPAddr) {
	msg.DestAddr = addr
}

func (msg *AcknowledgementMessage) GetDestination() net.UDPAddr {
	return msg.DestAddr
}

func (msg *AcknowledgementMessage) GetMessageType() int {
	return MESSAGE_ACKNOWLEDGE
}

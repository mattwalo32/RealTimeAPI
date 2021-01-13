package messages

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net"
)

const (
	MESSAGE_INVALID = iota
	MESSAGE_ACKNOWLEDGE
	MESSAGE_FIND_ROOM
	MESSAGE_JOIN_ROOM
	MESSAGE_LEAVE_ROOM
	MESSAGE_JOIN_SERVER
	MESSAGE_GAME_DATA
	MESSAGE_HEARTBEAT
	MESSAGE_STATUS
	MESSAGE_END_GAME
)

type Message interface {
	Encode() ([]byte, error)
	Decode([]byte)
	GetSource() net.UDPAddr
	SetSource(net.UDPAddr)
	GetDestination() net.UDPAddr
	SetDestination(net.UDPAddr)
	SetPacketNumber(int)
	GetPacketNumber() int
	IsResponseRequired() bool
	SetResponseRequired(bool)
	GetID() uuid.UUID
	GetMessageType() int
}

type ClientAssociable interface {
	GetClientID() uuid.UUID
}

type encodedMessage struct {
	MessageType int
	Data        []byte
}

func EncodeWithHeader(encodable Message) ([]byte, error) {
	data, err := encodable.Encode()
	if err != nil {
		return nil, err
	}

	message := encodedMessage{
		MessageType: encodable.GetMessageType(),
		Data:        data,
	}

	return json.Marshal(message)
}

func DecodeFromHeader(data []byte) (Message, error) {
	var header encodedMessage
	err := json.Unmarshal(data, &header)

	switch header.MessageType {
	case MESSAGE_FIND_ROOM:
		message := &FindRoomMessage{}
		message.Decode(header.Data)
		return message, err
	case MESSAGE_ACKNOWLEDGE:
		message := &AcknowledgementMessage{}
		message.Decode(header.Data)
		return message, err
	case MESSAGE_JOIN_SERVER:
		message := &JoinServerMessage{}
		message.Decode(header.Data)
		return message, err
	default:
		return nil, fmt.Errorf("Cannot decode unrecognized message type: %v", header.MessageType)
	}
}

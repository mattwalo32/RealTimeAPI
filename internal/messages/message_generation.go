package messages

import (
	"github.com/google/uuid"
	"github.com/mattwalo32/RealTimeAPI/internal/util"
	"math/rand"
)

func RandRoomMessage() *FindRoomMessage {
	return &FindRoomMessage{
		SourceAddr:       util.RandUDPAddr(),
		DestAddr:         util.RandUDPAddr(),
		MessageID:        uuid.New(),
		PacketNumber:     rand.Int(),
		ResponseRequired: util.RandBool(),

		UserID:              uuid.New(),
		ShouldStartWhenFull: util.RandBool(),
		MinPlayers:          rand.Int(),
		MaxPlayers:          rand.Int(),
	}
}

func RandAcknowledgeMessage() *AcknowledgementMessage {
	return &AcknowledgementMessage{
		SourceAddr:       util.RandUDPAddr(),
		DestAddr:         util.RandUDPAddr(),
		MessageID:        uuid.New(),
		PacketNumber:     rand.Int(),
		ResponseRequired: util.RandBool(),

		UserID:                uuid.New(),
		AcknowledgedMessageID: uuid.New(),
	}
}

func RandJoinServerMessage() *AcknowledgementMessage {
	return &AcknowledgementMessage{
		SourceAddr:       util.RandUDPAddr(),
		DestAddr:         util.RandUDPAddr(),
		MessageID:        uuid.New(),
		PacketNumber:     rand.Int(),
		ResponseRequired: util.RandBool(),
	}
}

func RandGameMessage() Encodable {
	messageType := rand.Intn(1)

	switch messageType {
	case 0:
		return RandRoomMessage()
	}

	return RandRoomMessage()
}

func RandEncodable() Encodable {
	messageType := rand.Intn(3)

	switch messageType {
	case 0:
		return RandRoomMessage()
	case 1:
		return RandAcknowledgeMessage()
	case 2:
		return RandJoinServerMessage()
	}

	return RandRoomMessage()
}

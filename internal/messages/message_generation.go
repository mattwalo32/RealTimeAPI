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

func RandEncodable() Encodable {
	// TODO: Do with all message types later
	return RandRoomMessage()
}

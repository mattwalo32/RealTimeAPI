package messages

import (
	"math/rand"
	"github.com/google/uuid"
	"github.com/mattwalo32/RealTimeAPI/internal/util"
)

func RandRoomMessage() *FindRoomMessage {
	return &FindRoomMessage{
		SourceAddr: util.RandUDPAddr(),
		DestAddr: util.RandUDPAddr(),
	
		UserID: uuid.New(),
		ShouldStartWhenFull: util.RandBool(),
		MinPlayers:          rand.Int(),
		MaxPlayers:          rand.Int(),
	} 
}

func RandEncodable() Encodable {
	// TODO: Do with all message types later
	return RandRoomMessage()
}
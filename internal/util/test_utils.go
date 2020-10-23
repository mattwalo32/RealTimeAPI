package util

import (
	"net"
	"math"
	"math/rand"
	"github.com/google/uuid"
	"github.com/mattwalo32/RealTimeAPI/internal/messages"
)

func generateRoomMessage() messages.FindRoomMessage {
	return messages.FindRoomMessage{
		SourceAddr: RandUDPAddr(),
		DestAddr: RandUDPAddr(),
	
		UserID: uuid.New(),
		ShouldStartWhenFull: RandBool(),
		MinPlayers:          rand.Int(),
		MaxPlayers:          rand.Int(),
	} 
}

func RandUDPAddr() net.UDPAddr {
	ip := []byte{byte(rand.Int()), byte(rand.Int()), byte(rand.Int()), byte(rand.Int())}

	return net.UDPAddr{
		IP: ip, 
		Port: rand.Int() % math.MaxUint16, 
		Zone: RandString(126),
	}
}

func RandString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(rand.Int())
	}
	
	return string(bytes)
}

func RandBool() bool {
    return rand.Intn(2) == 1
}
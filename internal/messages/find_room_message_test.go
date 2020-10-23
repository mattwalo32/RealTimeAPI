package messages

import (
	"net"
	"testing"
	"github.com/google/uuid"
)

var testAddressA = net.UDPAddr{[]byte{127, 0, 0, 1}, 50, ""}
var testAddressB = net.UDPAddr{[]byte{255, 1, 2, 28}, 0, "zone"}

func TestFindRoom_Type(t *testing.T) {
	msg := FindRoomMessage {}

	if msg.GetMessageType() != MESSAGE_FIND_ROOM {
		t.Errorf("Expected message type %v, got: %v", MESSAGE_FIND_ROOM, msg.GetMessageType())
	}
}

func TestFindRoom_EncodeDecode(t *testing.T) {
	msg := FindRoomMessage{
		SourceAddr: testAddressA,
		DestAddr: testAddressB,
		UserID: uuid.New(),
		ShouldStartWhenFull: true,
		MinPlayers: 1,
		MaxPlayers: 2,
	}

	data, err := msg.Encode()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	decodedMsg := FindRoomMessage{}
	decodedMsg.Decode(data)

	if (msg.SourceAddr.String() != decodedMsg.SourceAddr.String()) {
		t.Error("Source addresses do not match")
	}

	if (msg.DestAddr.String() != decodedMsg.DestAddr.String()) {
		t.Error("Destination addresses do not match")
	}

	if (msg.UserID != decodedMsg.UserID) {
		t.Error("UserIDs do not match")
	}

	if (msg.ShouldStartWhenFull != decodedMsg.ShouldStartWhenFull) {
		t.Error("ShouldStartWhenFull does not match")
	}

	if (msg.MinPlayers != decodedMsg.MinPlayers) {
		t.Error("MinPlayers does not match")
	}

	if (msg.MaxPlayers != decodedMsg.MaxPlayers) {
		t.Error("MaxPlayers does not match")
	}
}
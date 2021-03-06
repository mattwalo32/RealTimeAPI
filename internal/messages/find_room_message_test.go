package messages

import (
	"testing"
	"reflect"
)

func TestFindRoom_Type(t *testing.T) {
	msg := FindRoomMessage{}

	if msg.GetMessageType() != MESSAGE_FIND_ROOM {
		t.Errorf("Expected message type %v, got: %v", MESSAGE_FIND_ROOM, msg.GetMessageType())
	}
}

func TestFindRoom_Routable(t *testing.T) {
	msg := RandRoomMessage()

	func(routableMsg RoutableMessage) {
		if routableMsg.GetClientID() != msg.ClientID {
			t.Errorf("Message client ID and associated client ID do not match")
		}
	}(msg)
}

func TestFindRoom_EncodeDecode(t *testing.T) {
	msg := RandRoomMessage()

	data, err := msg.Encode()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	decodedMsg := FindRoomMessage{}
	decodedMsg.Decode(data)

	if msg.SourceAddr.String() != decodedMsg.SourceAddr.String() {
		t.Error("Source addresses do not match")
	}

	if msg.DestAddr.String() != decodedMsg.DestAddr.String() {
		t.Error("Destination addresses do not match")
	}

	if msg.MessageID != decodedMsg.MessageID {
		t.Error("MessageIDs do not match")
	}

	if msg.PacketNumber != decodedMsg.PacketNumber {
		t.Error("PacketNumbers do not match")
	}

	if msg.ResponseRequired != decodedMsg.ResponseRequired {
		t.Error("ResponseRequireds do not match")
	}

	if msg.UserID != decodedMsg.UserID {
		t.Error("UserIDs do not match")
	}

	if msg.ClientID != decodedMsg.ClientID {
		t.Error("ClientIDs do not match")
	}

	if msg.ShouldStartWhenFull != decodedMsg.ShouldStartWhenFull {
		t.Error("ShouldStartWhenFull does not match")
	}

	if !reflect.DeepEqual(msg.RoomTypes, decodedMsg.RoomTypes) {
		t.Error("Room Types do not match")
	}
}

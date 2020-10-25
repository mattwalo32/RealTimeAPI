package messages

import (
	"testing"
)

func TestJoinServer_Type(t *testing.T) {
	msg := JoinServerMessage{}

	if msg.GetMessageType() != MESSAGE_JOIN_SERVER {
		t.Errorf("Expected message type %v, got: %v", MESSAGE_JOIN_SERVER, msg.GetMessageType())
	}
}

func TestJoinServer_EncodeDecode(t *testing.T) {
	msg := RandJoinServerMessage()

	data, err := msg.Encode()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	decodedMsg := JoinServerMessage{}
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
}

package messages

import (
	"testing"
)

func TestAcknowledge_Type(t *testing.T) {
	msg := AcknowledgementMessage{}

	if msg.GetMessageType() != MESSAGE_ACKNOWLEDGE {
		t.Errorf("Expected message type %v, got: %v", MESSAGE_ACKNOWLEDGE, msg.GetMessageType())
	}
}

func TestAcknowledge_EncodeDecode(t *testing.T) {
	msg := RandAcknowledgeMessage()

	data, err := msg.Encode()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	decodedMsg := AcknowledgementMessage{}
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

	if msg.AcknowledgedMessageID != decodedMsg.AcknowledgedMessageID {
		t.Error("Acknowledged Message IDs do not match")
	}
}

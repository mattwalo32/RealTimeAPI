package server_test

import (
	"bytes"
	"github.com/mattwalo32/RealTimeAPI/internal/messages"
	"github.com/mattwalo32/RealTimeAPI/pkg/server"
	"net"
	"testing"
)

func createMessageHandler(address string) (chan messages.Encodable, *server.MessageHandler) {
	receivingChan := make(chan messages.Encodable, 2)
	config := server.MessageHandlerConfig{
		MessageReceivingChan: receivingChan,
		Address:              address,
	}

	return receivingChan, server.NewMessageHandler(config)
}

func TestSendMessages(t *testing.T) {
	clientAAddress := "localhost:9999"
	clientBAddress := "localhost:9998"
	numTestMessages := 20
	_, handlerA := createMessageHandler(clientAAddress)
	clientBReceivingChan, _ := createMessageHandler(clientBAddress)

	clientAUDPAddr, _ := net.ResolveUDPAddr("udp4", clientAAddress)
	clientBUDPAddr, _ := net.ResolveUDPAddr("udp4", clientBAddress)

	test_messages := make([]messages.Encodable, numTestMessages)
	for i := 0; i < 20; i++ {
		test_messages[i] = messages.RandEncodable()
		test_messages[i].SetSource(*clientAUDPAddr)
		test_messages[i].SetDestination(*clientBUDPAddr)
	}

	for packetNum, msg := range test_messages {
		handlerA.SendMessage(msg)
		response := <-clientBReceivingChan

		if response.GetMessageType() != msg.GetMessageType() {
			t.Fatalf("Expected message type %v, got: %v", msg.GetMessageType(), response.GetMessageType())
		}

		if response.GetPacketCount() != packetNum {
			t.Errorf("Expected packet number of %v, got: %v", packetNum, response.GetPacketCount())
		}

		source := response.GetSource()
		if source.String() != clientAUDPAddr.String() {
			t.Errorf("Expected message to be from %v, got: %v", clientAUDPAddr, response.GetSource())
		}

		expectedBytes, _ := msg.Encode()
		actualBytes, _ := response.Encode()
		if !bytes.Equal(expectedBytes, actualBytes) {
			t.Errorf("The content of the sent and received messages differ")
		}
	}
}

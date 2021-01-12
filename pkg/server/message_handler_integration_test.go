package server

import (
	"bytes"
	"github.com/mattwalo32/RealTimeAPI/internal/messages"
	"net"
	"testing"
	"time"
)

const (
	RESPONSE_WAIT_TIME = 1000 * time.Millisecond
)

func createMessageHandler(address string) (chan messages.Message, *MessageHandler) {
	receivingChan := make(chan messages.Message, 10)
	config := MessageHandlerConfig{
		MessageReceivingChan:  receivingChan,
		MaxMessageRetries:     5,
		MessageRetryTimeoutMs: uint64(500),
		Address:               address,
	}

	return receivingChan, NewMessageHandler(config)
}

/**
  * Test sending a game-related message unreliably
  */
func TestSendGameMessages_Unreliable(t *testing.T) {
	clientAAddress := "localhost:9999"
	clientBAddress := "localhost:9998"
	numTestMessages := 20
	_, handlerA := createMessageHandler(clientAAddress)
	clientBReceivingChan, handlerB := createMessageHandler(clientBAddress)

	clientAUDPAddr, _ := net.ResolveUDPAddr("udp4", clientAAddress)
	clientBUDPAddr, _ := net.ResolveUDPAddr("udp4", clientBAddress)

	test_messages := make([]messages.Message, numTestMessages)
	for i := 0; i < numTestMessages; i++ {
		test_messages[i] = messages.RandGameMessage()
		test_messages[i].SetSource(*clientAUDPAddr)
		test_messages[i].SetDestination(*clientBUDPAddr)
	}

	for packetNum, msg := range test_messages {
		handlerA.SendMessageUnreliably(msg)
		response := <-clientBReceivingChan

		if response.GetMessageType() != msg.GetMessageType() {
			t.Fatalf("Expected message type %v, got: %v", msg.GetMessageType(), response.GetMessageType())
		}

		if response.GetPacketNumber() != packetNum {
			t.Errorf("Expected packet number of %v, got: %v", packetNum, response.GetPacketNumber())
		}

		if response.IsResponseRequired() {
			t.Errorf("Expected response not required")
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

	handlerA.Stop()
	handlerB.Stop()
}

// TODO: Test send Non-Game Messages

/**
  * Test sending a game message reliably with no response. The message will
  * be resent multiple times until we reach the retry limit.
  */
func TestSendGameMessages_Reliable_NoResponse(t *testing.T) {
	clientAAddress := "localhost:9999"
	clientBAddress := "localhost:9998"
	numTestMessages := 20
	_, handlerA := createMessageHandler(clientAAddress)

	clientAUDPAddr, _ := net.ResolveUDPAddr("udp4", clientAAddress)
	clientBUDPAddr, _ := net.ResolveUDPAddr("udp4", clientBAddress)

	test_messages := make([]messages.Message, numTestMessages)
	for i := 0; i < numTestMessages; i++ {
		test_messages[i] = messages.RandGameMessage()
		test_messages[i].SetSource(*clientAUDPAddr)
		test_messages[i].SetDestination(*clientBUDPAddr)
	}

	for _, msg := range test_messages {
		handlerA.SendMessageReliably(msg)
	}

	if len(handlerA.messageRetryEventIDs) != numTestMessages {
		t.Errorf("Expected %v retry event IDs", numTestMessages)
	}

	if handlerA.timer.NumEvents() != numTestMessages {
		t.Errorf("Expected %v timer events", numTestMessages)
	}

	handlerA.Stop()
}

/**
  * Test sending a game message reliably with a response. It should only be sent once.
  */
func TestSendGameMessages_Reliable_Response(t *testing.T) {
	clientAAddress := "localhost:9999"
	clientBAddress := "localhost:9998"
	numTestMessages := 3
	_, handlerA := createMessageHandler(clientAAddress)
	_, handlerB := createMessageHandler(clientBAddress)

	clientAUDPAddr, _ := net.ResolveUDPAddr("udp4", clientAAddress)
	clientBUDPAddr, _ := net.ResolveUDPAddr("udp4", clientBAddress)

	test_messages := make([]messages.Message, numTestMessages)
	for i := 0; i < numTestMessages; i++ {
		test_messages[i] = messages.RandMessage()
		test_messages[i].SetSource(*clientAUDPAddr)
		test_messages[i].SetDestination(*clientBUDPAddr)
	}

	for _, msg := range test_messages {
		handlerA.SendMessageReliably(msg)
	}

	<-time.After(RESPONSE_WAIT_TIME)

	if len(handlerA.messageRetryEventIDs) != 0 {
		t.Errorf("%v event IDs were not deleted", len(handlerA.messageRetryEventIDs))
	}

	if handlerA.timer.NumEvents() != 0 {
		t.Errorf("%v events were not deleted from timer", handlerA.timer.NumEvents())
	}

	handlerA.Stop()
	handlerB.Stop()
}
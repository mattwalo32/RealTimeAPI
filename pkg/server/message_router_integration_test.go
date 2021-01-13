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

func createMessageRouter(address string) (chan messages.Message, *MessageRouter) {
	receivingChan := make(chan messages.Message, 10)
	config := MessageRouterConfig{
		MaxMessageRetries:     5,
		MessageRetryTimeoutMs: uint64(500),
		Address:               address,
	}

	return receivingChan, NewMessageRouter(config)
}

/**
  * Test sending a game-related message unreliably
  */
func TestSendGameMessages_Unreliable(t *testing.T) {
	clientAAddress := "localhost:9999"
	clientBAddress := "localhost:9998"
	numTestMessages := 20
	_, routerA := createMessageRouter(clientAAddress)
	clientBReceivingChan, routerB := createMessageRouter(clientBAddress)

	clientAUDPAddr, _ := net.ResolveUDPAddr("udp4", clientAAddress)
	clientBUDPAddr, _ := net.ResolveUDPAddr("udp4", clientBAddress)

	test_messages := make([]messages.Message, numTestMessages)
	for i := 0; i < numTestMessages; i++ {
		test_messages[i] = messages.RandGameMessage()
		test_messages[i].SetSource(*clientAUDPAddr)
		test_messages[i].SetDestination(*clientBUDPAddr)
	}

	for packetNum, msg := range test_messages {
		routerA.SendMessageUnreliably(msg)
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

	routerA.Stop()
	routerB.Stop()
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
	_, routerA := createMessageRouter(clientAAddress)

	clientAUDPAddr, _ := net.ResolveUDPAddr("udp4", clientAAddress)
	clientBUDPAddr, _ := net.ResolveUDPAddr("udp4", clientBAddress)

	test_messages := make([]messages.Message, numTestMessages)
	for i := 0; i < numTestMessages; i++ {
		test_messages[i] = messages.RandGameMessage()
		test_messages[i].SetSource(*clientAUDPAddr)
		test_messages[i].SetDestination(*clientBUDPAddr)
	}

	for _, msg := range test_messages {
		routerA.SendMessageReliably(msg)
	}

	if len(routerA.messageRetryEventIDs) != numTestMessages {
		t.Errorf("Expected %v retry event IDs", numTestMessages)
	}

	if routerA.timer.NumEvents() != numTestMessages {
		t.Errorf("Expected %v timer events", numTestMessages)
	}

	routerA.Stop()
}

/**
  * Test sending a game message reliably with a response. It should only be sent once.
  */
func TestSendGameMessages_Reliable_Response(t *testing.T) {
	clientAAddress := "localhost:9999"
	clientBAddress := "localhost:9998"
	numTestMessages := 3
	_, routerA := createMessageRouter(clientAAddress)
	_, routerB := createMessageRouter(clientBAddress)

	clientAUDPAddr, _ := net.ResolveUDPAddr("udp4", clientAAddress)
	clientBUDPAddr, _ := net.ResolveUDPAddr("udp4", clientBAddress)

	test_messages := make([]messages.Message, numTestMessages)
	for i := 0; i < numTestMessages; i++ {
		test_messages[i] = messages.RandMessage()
		test_messages[i].SetSource(*clientAUDPAddr)
		test_messages[i].SetDestination(*clientBUDPAddr)
	}

	for _, msg := range test_messages {
		routerA.SendMessageReliably(msg)
	}

	<-time.After(RESPONSE_WAIT_TIME)

	if len(routerA.messageRetryEventIDs) != 0 {
		t.Errorf("%v event IDs were not deleted", len(routerA.messageRetryEventIDs))
	}

	if routerA.timer.NumEvents() != 0 {
		t.Errorf("%v events were not deleted from timer", routerA.timer.NumEvents())
	}

	routerA.Stop()
	routerB.Stop()
}
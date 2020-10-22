package conn_test

import (
	"bytes"
	"github.com/mattwalo32/RealTimeAPI/internal/conn"
	"net"
	"testing"
)

var ()

func createUDPManager(address string) (chan conn.Message, *conn.UDPManager) {
	receivingChan := make(chan conn.Message, 2)
	config := conn.UDPManagerConfig{
		ReceivingChan: receivingChan,
		Address:       address,
	}

	return receivingChan, conn.NewUDPManager(config)
}

func TestSendMessages(t *testing.T) {
	clientAAddress := "localhost:9999"
	clientBAddress := "localhost:9998"

	_, managerA := createUDPManager(clientAAddress)
	clientBReceivingChan, _ := createUDPManager(clientBAddress)

	clientAUDPAddr, _ := net.ResolveUDPAddr("udp4", clientAAddress)
	clientBUDPAddr, _ := net.ResolveUDPAddr("udp4", clientBAddress)

	message_tests := [][]byte{
		[]byte(""),
		[]byte("Test"),
		[]byte("Longer message with more substance!"),
		[]byte("!@#$%&*()39283   add\t\tflkj \n !!/?"),
	}

	for _, testData := range message_tests {
		msg := conn.Message{
			Data:    testData,
			Address: *clientBUDPAddr,
		}

		managerA.SendMessage(msg)
		response := <-clientBReceivingChan

		if !bytes.Equal(response.Data, testData) {
			t.Errorf("Expected message data (%v) is different from actual (%v)", testData, response.Data)
		}

		if response.Address.String() != clientAUDPAddr.String() {
			t.Errorf("Expected message to be from %v, got: %v", clientAAddress, response.Address)
		}
	}
}

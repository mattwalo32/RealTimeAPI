package conn

import (
	"testing"
)

func TestAssertConfigValid_TooShort(t *testing.T) {
	config := &UDPManagerConfig{
		ReceivingChan: make(chan Message, MIN_RECEIVING_CHAN_CAP - 1),
		Address: "",
	}

	err := assertConfigValid(config)
	if err == nil {
		t.Errorf("Expected error, receiving chan capacity is too low")
	}
}

func TestAssertConfigValid_Valid(t *testing.T) {
	config := &UDPManagerConfig{
		ReceivingChan: make(chan Message, MIN_RECEIVING_CHAN_CAP),
		Address: "",
	}

	err := assertConfigValid(config)
	if err != nil {
		t.Errorf("Error not expected")
	}
}
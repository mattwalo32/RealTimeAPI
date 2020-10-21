package main

import (
	"os"
	"net"
	"log"
	"fmt"
	"errors"
	"github.com/mattwalo32/RealTimeAPI/internal/conn"
)

const (
	MIN_NUMBER_ARGS = 3
)

var (
	doneChan = make(chan bool)
	errorMissingArgs = errors.New("Please provide the following args in order: <listen address> <write address>")
)

func main() {
	receivingChan := make(chan conn.Message)
	manager := createUDPManager(receivingChan)

	go printIncomingMessages(receivingChan)
	go sendUserInput(manager)
}

func createUDPManager(receivingChan chan conn.Message) *conn.AbstractUDPManager {
	if len(os.Args) < MIN_NUMBER_ARGS + 1 {
		log.Fatal(errorMissingArgs)
	}

	address := os.Args[1]
	config := conn.UDPManagerConfig {
		ReceivingChan: receivingChan,
		Address: address,
	}

	return conn.NewUDPManager(config)
}

func printIncomingMessages(receivingChan chan conn.Message) {
	for {
		select {
			case <- doneChan:
				return
			case msg := <- receivingChan:
				fmt.Println(msg.Data)
		}
	}
}

func sendUserInput(manager *conn.AbstractUDPManager) {
	writeAddress := resolveUDPAddr(os.Args[2])
	var input string

	for {
		fmt.Scanln(&input)

		msg := conn.Message{
			Data: []byte(input),
			Address: writeAddress,
		}

		manager.SendMessage(msg)
	}
}

func resolveUDPAddr(address string) net.UDPAddr {
	addr, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		log.Fatalf("Could not resolve address: %v", err)
	}

	return *addr
}
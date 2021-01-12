package main

import (
	"errors"
	"fmt"
	"github.com/mattwalo32/RealTimeAPI/internal/conn"
	"log"
	"net"
	"os"
	"os/signal"
)

const (
	MIN_NUMBER_ARGS = 2
)

var (
	doneChan         = make(chan bool)
	errorMissingArgs = errors.New("Please provide the following args in order: <listen address> <write address>")
)

func main() {
	receivingChan := make(chan conn.Packet, 2)
	manager := createUDPManager(receivingChan)

	go printIncomingPackets(receivingChan)
	go sendUserInput(manager)

	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt)

	<-interrupt
	manager.Close()
	close(doneChan)
}

func createUDPManager(receivingChan chan conn.Packet) *conn.UDPManager {
	if len(os.Args) < MIN_NUMBER_ARGS+1 {
		log.Fatal(errorMissingArgs)
	}

	address := os.Args[1]
	config := conn.UDPManagerConfig{
		ReceivingChan: receivingChan,
		Address:       address,
	}

	return conn.NewUDPManager(config)
}

func printIncomingPackets(receivingChan chan conn.Packet) {
	for {
		select {
		case <-doneChan:
			return
		case msg := <-receivingChan:
			fmt.Println(string(msg.Data))
		}
	}
}

func sendUserInput(manager *conn.UDPManager) {
	writeAddress := resolveUDPAddr(os.Args[2])
	var input string

	for {
		fmt.Scanln(&input)

		msg := conn.Packet{
			Data:    []byte(input),
			Address: writeAddress,
		}

		manager.SendPacket(msg)
	}
}

func resolveUDPAddr(address string) net.UDPAddr {
	addr, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		log.Fatalf("Could not resolve address: %v", err)
	}

	return *addr
}

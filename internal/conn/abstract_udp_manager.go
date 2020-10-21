package conn

import (
	"fmt"
	"net"
	"log"
)

const (
	MIN_RECEIVING_CHAN_CAP = 2
	SENDING_CHAN_CAP = 2
	BUFFER_SIZE = 1024

)

type AbstractUDPManager struct {
	config *UDPManagerConfig
	conn   *net.UDPConn

	sendingChan chan Message
	doneChan chan bool
}

type UDPManagerConfig struct {
	// The manager writes incoming messages to this channel. Acts as callback mechanism.
	ReceivingChan chan Message

	Address string
}

type Message struct {
	Data []byte
	Address net.UDPAddr
}

func NewUDPManager(config UDPManagerConfig) *AbstractUDPManager {
	err := assertConfigValid(&config)
	if (err != nil) {
		log.Fatal(err)
	}
	
	manager := &AbstractUDPManager{
		config: &config,
		sendingChan: make(chan Message, SENDING_CHAN_CAP),
		doneChan: make(chan bool),
	}

	manager.init()

	return manager
}

func assertConfigValid(config *UDPManagerConfig) error {
	cap := cap(config.ReceivingChan)
	if (cap < MIN_RECEIVING_CHAN_CAP) {
		return fmt.Errorf("Recieving channel must have capacity %v, got: %v", MIN_RECEIVING_CHAN_CAP, cap)
	}

	return nil
}

func (manager *AbstractUDPManager) init() {
	manager.initConn()
	go manager.listenUDP()
	go manager.sendUDP()
}

func (manager *AbstractUDPManager) initConn() {
	address := manager.config.Address
	addr, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error resolving address (%v): %v", address, err))
	}

	manager.conn, err = net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error listening on UDP: %v", err))
	}
}

func (manager *AbstractUDPManager) listenUDP() {
	buffer := make([]byte, BUFFER_SIZE)

	for {
		select {
			case <-manager.doneChan:
				return
			default:
		}

		n, addr, err := manager.conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("WARNING: Error while reading from UDP: %v", err)
			return
		}

		msg :=  Message {
			Data: buffer[0:n], 
			Address: *addr,
		}

		manager.config.ReceivingChan <- msg
	}
	// TODO: Put messages into queue if sendingChan is full
}

func queueMessageForProcessing() {
	// TODO:
}

func (manager *AbstractUDPManager) sendUDP() {
	msg := Message{}

	for {
		select {
			case <-manager.doneChan:
				return
			case msg = <-manager.sendingChan:
		}

		_,err := manager.conn.WriteToUDP(msg.Data, &msg.Address)
		if err != nil {
			log.Printf("WARNING: Error while writing to UDP: %v", err)
		}
	}
}

func (manager *AbstractUDPManager) SendMessage(msg Message) {
	manager.sendingChan <- msg
}

func (manager *AbstractUDPManager) Close() {
	close(manager.doneChan)
	manager.conn.Close()
}
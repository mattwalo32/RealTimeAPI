package conn

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
)

const (
	MIN_RECEIVING_CHAN_CAP = 2
	SENDING_CHAN_CAP       = 15
	BUFFER_SIZE            = 1024
)

type UDPManager struct {
	config *UDPManagerConfig
	conn   *net.UDPConn
	addr *net.UDPAddr

	sendingChan chan Packet
	doneChan    chan bool
}

type UDPManagerConfig struct {
	// The manager writes incoming messages to this channel. Acts as callback mechanism.
	ReceivingChan chan Packet

	// Address to listen on
	Address string
}

type Packet struct {
	Data    []byte
	Address net.UDPAddr
}

func NewUDPManager(config UDPManagerConfig) *UDPManager {
	err := assertConfigValid(&config)
	if err != nil {
		log.Fatal(err)
	}

	manager := &UDPManager{
		config:      &config,
		sendingChan: make(chan Packet, SENDING_CHAN_CAP),
		doneChan:    make(chan bool),
	}

	manager.init()

	return manager
}

func assertConfigValid(config *UDPManagerConfig) error {
	cap := cap(config.ReceivingChan)
	if cap < MIN_RECEIVING_CHAN_CAP {
		return fmt.Errorf("Recieving channel must have capacity %v, got: %v", MIN_RECEIVING_CHAN_CAP, cap)
	}

	return nil
}

func (manager *UDPManager) init() {
	manager.initConn()
	go manager.listenUDP()
	go manager.sendUDP()
}

func (manager *UDPManager) initConn() {
	address := manager.config.Address
	addr, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error resolving address (%v): %v", address, err))
	}

	manager.addr = addr
	manager.conn, err = net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error listening on UDP: %v", err))
	}
}

func (manager *UDPManager) listenUDP() {
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

		bufferCopy := make([]byte, n)
		copy(bufferCopy, buffer[0:n])

		msg := Packet{
			Data:    bufferCopy,
			Address: *addr,
		}

		log.WithFields(log.Fields{
			"Data": msg.Data,
			"Address": msg.Address,
		}).Trace("Received message via UDP")

		manager.config.ReceivingChan <- msg
	}
}

func (manager *UDPManager) sendUDP() {
	msg := Packet{}

	for {
		select {
		case <-manager.doneChan:
			return
		case msg = <-manager.sendingChan:
		}

		log.WithFields(log.Fields{
			"Data": msg.Data,
			"Address": msg.Address,
		}).Trace("Sending message via UDP")

		_, err := manager.conn.WriteToUDP(msg.Data, &msg.Address)
		if err != nil {
			log.Printf("WARNING: Error while writing to UDP: %v", err)
		}
	}
}

func (manager *UDPManager) SendPacket(msg Packet) {
	manager.sendingChan <- msg
}

func (manager *UDPManager) GetUDPAddr() *net.UDPAddr {
	return manager.addr
}

func (manager *UDPManager) Stop() {
	close(manager.doneChan)
	manager.conn.Close()
}

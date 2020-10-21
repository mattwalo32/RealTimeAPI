package server

import (
	"net"
)

type RealTimeServer struct {
	// Configuration passed into server
	config *RealTimeServerConfig

	// TODO: Abstract UDP Manager

	// Close channel to stop all goroutines
	doneChan chan bool
}

type RealTimeServerConfig struct {
	UDPListeningAddress string
	TCPListeningAddress string
}

func NewServer(config RealTimeServerConfig) *RealTimeServer {
	server := &RealTimeServer {
		config: &config,
		doneChan: make(chan bool),
	}

	return server
}

func (r *RealTimeServer) Start() {
	// Abstract away UDP/TCP connections
}

func (r *RealTimeServer) Stop() {
	close(r.doneChan)
}
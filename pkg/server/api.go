package server

type RealTimeServer struct {
	doneChan chan bool
}

type RealTimeServerConfig struct {

}

func NewServer() *RealTimeServer {
	server := &RealTimeServer {
		doneChan: make(chan bool),
	}

	return server
}

func Start() {

}

func Stop() {

}
package sse

import "log"

type SseServer struct {
	// channel of messages to be pushed to clients
	messages chan []byte

	// map of all connected clients
	// - key is the channel of bytes which handles the
	//   comnunication between the SseServer and the client
	// - value is meaningless
	clients map[chan []byte]bool

	// channel to regiter new client
	register chan chan []byte

	// channel to unregister clients
	unregister chan chan []byte
}

// NewServer create and start a new SseServer
func NewSseServer() *SseServer {
	server := &SseServer{
		messages:   make(chan []byte, 1),
		clients:    make(map[chan []byte]bool),
		register:   make(chan chan []byte),
		unregister: make(chan chan []byte),
	}
	go server.listen()
	return server
}

func (server *SseServer) listen() {
	for {
		select {
		case c := <-server.messages:
			log.Println("new message")
		case c := <-server.register:
			log.Println("new client")
		case c := <-server.unregister:
			log.Println("client leaving")
		}
	}

}

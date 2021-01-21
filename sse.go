package sse

import (
	"log"
	"net/http"
)

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

// Publish message to all connected clients
func (server *SseServer) Publish(message string) {
	server.messages <- []byte(message)
}
func (server *SseServer) Connect(rw http.ResponseWriter, req *http.Request) {
	flusher, ok := rw.(http.Flusher)

	if !ok {
		http.Error(rw, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	// Set the headers related to event streaming.
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")

	// create communication channel and register
	c := make(chan []byte)
	server.register <- c

	// listen to the closing of the http connection via the CloseNotifier
	notify := rw.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify
		log.Println("connection closed")
		server.unregister <- c
	}()

	for {
		_, err := rw.Write(<-c)
		if err != nil {
			log.Printf("%v", err)
		}
		flusher.Flush()
	}
}
func (server *SseServer) broadcast(message []byte) {
	for c := range server.clients {
		c <- message
	}
}
func (server *SseServer) listen() {
	for {
		select {
		case c := <-server.messages:
			log.Println("new message")
			server.broadcast(c)
		case c := <-server.register:
			log.Println("new client")
			server.clients[c] = true
		case c := <-server.unregister:
			log.Println("client leaving")
			delete(server.clients, c)
		}
	}
}

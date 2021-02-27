package sse

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Server ...
type Server struct {
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
func NewServer() *Server {
	server := &Server{
		messages:   make(chan []byte, 1),
		clients:    make(map[chan []byte]bool),
		register:   make(chan chan []byte),
		unregister: make(chan chan []byte),
	}
	go server.listen()
	return server
}

// Publish publish object. mesage will be conform to
// EventSource spec:
// data:[json]\n\n
func (s *Server) Publish(message interface{}) error {
	var sb strings.Builder
	sb.WriteString("data:")
	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	sb.Write(bytes)
	sb.WriteString("\n\n")

	s.messages <- []byte(sb.String())
	return nil
}

// PublishRaw publish raw message to all connected clients
func (s *Server) PublishRaw(message string) {
	s.messages <- []byte(message)
}

// Subscribe to the sse server
func (s *Server) Subscribe(rw http.ResponseWriter, req *http.Request) {
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
	s.register <- c

	// listen to the closing of the http connection via the CloseNotifier
	notify := rw.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify
		log.Println("connection closed")
		s.unregister <- c
	}()

	for msg := range c {
		_, err := fmt.Fprintf(rw, "%s\n", msg)
		if err != nil {
			log.Printf("%v", err)
			s.unregister <- c
			break
		}
		flusher.Flush()
	}
}
func (s *Server) broadcast(message []byte) {
	for c := range s.clients {
		c <- message
	}
}
func (s *Server) listen() {
	for {
		select {
		case c := <-s.messages:
			// log.Println("new message")
			s.broadcast(c)
		case c := <-s.register:
			log.Println("client registering")
			s.clients[c] = true
		case c := <-s.unregister:
			log.Println("client leaving")
			delete(s.clients, c)
			close(c)
		}
	}
}

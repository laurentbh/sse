package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/laurentbh/sse"
)

func main() {

	server := sse.NewServer()

	http.HandleFunc("/time", server.Subscribe)

	go func() {
		ticker := time.Tick(1 * time.Second)
		count := 0
		for {
			select {
			case <-ticker:
				server.Publish(fmt.Sprintf("message # %d", count))
				count++
			}
		}
	}()

	http.ListenAndServe(":8080", http.DefaultServeMux)
}

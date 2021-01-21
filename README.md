# sse

Go implementation of SSE server


# install
`go get github.com/laurentbh/sse`

# Usage
1- start the server
```go
server := sse.NewSseServer()
```
2- connect the server to the endpoint
```go
http.HandleFunc("/time", server.Subscribe)
```
3- push message to the server
```go
server.Publish("message")
```

check [server example](example/sse_time.go)

with `curl -X GET http://localhost:8080/time` as client

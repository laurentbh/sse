# sse

Go implementation of SSE server


# install
`go get github.com/laurentbh/sse`

# Usage
1- start the server
```go
server := sse.NewServer()
```
2- connect the server to the endpoint
```go
http.HandleFunc("/time", server.Subscribe)
```
3- push message to the server

the `Publish` method will push messages compliant with `EventSource`, ie the message starts with the string `data:`, followed by json representation of the object, followed by 2 `\n`

```go
type test struct {
  ID   int
  Name string
}

t := test{ID: 1, Name: "test"}

server.Publish(t)
```
the `PublishRaw` method will push messages as passed without additional transformation

check [server example](example/sse_time.go)

with `curl -X GET http://localhost:8080/time` as client

# socket
[![Build Status](https://travis-ci.org/joaosoft/socket.svg?branch=master)](https://travis-ci.org/joaosoft/socket) | [![codecov](https://codecov.io/gh/joaosoft/socket/branch/master/graph/badge.svg)](https://codecov.io/gh/joaosoft/socket) | [![Go Report Card](https://goreportcard.com/badge/github.com/joaosoft/socket)](https://goreportcard.com/report/github.com/joaosoft/socket) | [![GoDoc](https://godoc.org/github.com/joaosoft/socket?status.svg)](https://godoc.org/github.com/joaosoft/socket)

A service that allows you to broadcast Publish and Listen messages to/from the server.

###### If i miss something or you have something interesting, please be part of this project. Let me know! My contact is at the end.

## How it works ?
> Flow example:
* CLIENT: creates a gateway server
* CLIENT: subscribe a topic/channel on the server 
* SERVER: register the client on the topic/channel with the gateway to call it back later
* CLIENT: send a message to a topic/channel
* SERVER: receives the message and broadcast the message to all registered clients
* CLIENT: all clients receive the message sent
* CLIENT: unsubscribe the topic/channel on the server 

## Usage 
This examples are available in the project at [socket/examples](https://github.com/joaosoft/socket/tree/master/examples)

### Server
```go
func main() {
	// server
	server, err := socket.NewServer()
	if err != nil {
		panic(err)
	}

	if err := server.Start(); err != nil {
		panic(err)
	}
}
```

### Client
```go
func main() {
	client, err := socket.NewClient()
	if err != nil {
		panic(err)
	}

	if err := client.Start(); err != nil {
		panic(err)
	}

	if err := client.Subscribe("topic_1", "channel_1"); err != nil {
		panic(err)
	}

	client.Listen("topic_1", "channel_1", func(message []byte) error {
		fmt.Printf("\nreceived on listener the message %s", string(message))
		return nil
	})

	if err := client.Publish("topic_1", "channel_1", []byte("hello, this is a test message")); err != nil {
		panic(err)
	}

	// Unsubscribe
	//if err := client.Unsubscribe("topic_1", "channel_1"); err != nil {
	//	panic(err)
	//}

	client.Wait()
}
```

## Dependecy Management
>### Dependency

Project dependencies are managed using Dep. Read more about [Dep](https://github.com/golang/dep).
* Get dependency manager: `go get github.com/joaosoft/dependency`
* Install dependencies: `dependency get`


>### Go
```
go get github.com/joaosoft/socket
```

## Known issues

## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com

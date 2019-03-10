package main

import (
	"socket"
)

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

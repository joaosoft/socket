package main

import (
	"fmt"
	"socket"
)

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

	if err := client.Unsubscribe("topic_1", "channel_1"); err != nil {
		panic(err)
	}
}

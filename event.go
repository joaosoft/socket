package socket

type eventHandler func(topic, channel string, message []byte) error

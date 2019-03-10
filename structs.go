package socket

type MessageHandler func(message []byte) error

type messageAcknowledge struct {
	Acknowledge bool  `json:"acknowledge"`
	Error       error `json:"error,omitempty"`
}

type listener struct {
	gateway string
}

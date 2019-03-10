package socket

type MessageHandler func(message []byte) error

type messageAcknowledge struct {
	Acknowledge bool    `json:"acknowledge"`
	Errors      []error `json:"errors,omitempty"`
}

type listener struct {
	gateway string
}

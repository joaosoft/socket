package socket

import (
	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

// SocketClientOption ...
type SocketClientOption func(client *Client)

// Reconfigure ...
func (s *Client) Reconfigure(options ...SocketClientOption) {
	for _, option := range options {
		option(s)
	}
}

// WithClientConfiguration ...
func WithClientConfiguration(config *SocketConfig) SocketClientOption {
	return func(client *Client) {
		client.config = config.Client
	}
}

// WithClientLogger ...
func WithClientLogger(logger logger.ILogger) SocketClientOption {
	return func(client *Client) {
		client.logger = logger
		client.isLogExternal = true
	}
}

// WithClientLogLevel ...
func WithClientLogLevel(level logger.Level) SocketClientOption {
	return func(client *Client) {
		client.logger.SetLevel(level)
	}
}

// WithClientManager ...
func WithClientManager(mgr *manager.Manager) SocketClientOption {
	return func(client *Client) {
		client.pm = mgr
	}
}

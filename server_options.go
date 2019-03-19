package socket

import (
	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

// SocketServerOption ...
type SocketServerOption func(server *Server)

// Reconfigure ...
func (s *Server) Reconfigure(options ...SocketServerOption) {
	for _, option := range options {
		option(s)
	}
}

// WithServerConfiguration ...
func WithServerConfiguration(config *ServerConfig) SocketServerOption {
	return func(server *Server) {
		server.config = config
	}
}

// WithServerLogger ...
func WithServerLogger(logger logger.ILogger) SocketServerOption {
	return func(server *Server) {
		server.logger = logger
		server.isLogExternal = true
	}
}

// WithServerLogLevel ...
func WithServerLogLevel(level logger.Level) SocketServerOption {
	return func(server *Server) {
		server.logger.SetLevel(level)
	}
}

// WithServerManager ...
func WithServerManager(mgr *manager.Manager) SocketServerOption {
	return func(server *Server) {
		server.pm = mgr
	}
}

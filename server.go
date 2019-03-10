package socket

import (
	"sync"

	"github.com/joaosoft/web"

	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

type Server struct {
	listeners map[string]map[string]map[string]*listener
	client    *web.Client
	server    *web.Server

	config        *ServerConfig
	isLogExternal bool
	pm            *manager.Manager
	logger        logger.ILogger
	mux           sync.Mutex
}

// NewServer ...
func NewServer(options ...SocketServerOption) (*Server, error) {
	config, simpleConfig, err := NewConfig()
	client, err := web.NewClient()
	if err != nil {
		return nil, err
	}

	service := &Server{
		listeners: make(map[string]map[string]map[string]*listener),
		client:    client,

		pm:     manager.NewManager(manager.WithRunInBackground(false)),
		logger: logger.NewLogDefault("socket-server", logger.WarnLevel),
		config: config.Socket.Server,
	}

	if service.isLogExternal {
		service.pm.Reconfigure(manager.WithLogger(logger.Instance))
	}

	if err != nil {
		service.logger.Error(err.Error())
	} else if config.Socket.Server != nil {
		service.pm.AddConfig("config_app", simpleConfig)
		level, _ := logger.ParseLevel(config.Socket.Server.Log.Level)
		service.logger.Debugf("setting log level to %s", level)
		service.logger.Reconfigure(logger.WithLevel(level))
	} else {
		config.Socket.Server = &ServerConfig{}
	}

	service.Reconfigure(options...)

	webServer := service.pm.NewSimpleWebServer(service.config.Address)
	service.server = webServer.GetClient().(*web.Server)
	service.initController()

	service.pm.AddWeb("api_web", webServer)

	return service, nil
}

// Start ...
func (s *Server) Start() error {
	return s.pm.Start()
}

// Stop ...
func (s *Server) Stop() error {
	return s.pm.Stop()
}

package socket

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/joaosoft/web"

	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

type Client struct {
	listeners map[string]map[string]MessageHandler
	client    *web.Client
	server    *web.Server

	config        *ClientConfig
	isLogExternal bool
	pm            *manager.Manager
	logger        logger.ILogger
	mux           sync.Mutex
	started       bool
}

// NewClient ...
func NewClient(options ...SocketClientOption) (*Client, error) {
	config, simpleConfig, err := NewConfig()
	client, err := web.NewClient()
	if err != nil {
		return nil, err
	}

	service := &Client{
		listeners: make(map[string]map[string]MessageHandler),
		client:    client,

		pm:     manager.NewManager(manager.WithRunInBackground(true)),
		logger: logger.NewLogDefault("socket-client", logger.WarnLevel),
		config: config.Socket.Client,
	}

	if service.isLogExternal {
		service.pm.Reconfigure(manager.WithLogger(logger.Instance))
	}

	if err != nil {
		service.logger.Error(err.Error())
	} else if config.Socket.Client != nil {
		service.pm.AddConfig("config_app", simpleConfig)
		level, _ := logger.ParseLevel(config.Socket.Client.Log.Level)
		service.logger.Debugf("setting log level to %s", level)
		service.logger.Reconfigure(logger.WithLevel(level))
	} else {
		config.Socket.Client = &ClientConfig{}
	}

	service.Reconfigure(options...)

	port, err := web.GetFreePort()
	if err != nil {
		return nil, err
	}

	webServer := service.pm.NewSimpleWebServer(fmt.Sprintf(":%d", port))
	service.server = webServer.GetClient().(*web.Server)
	service.initController()

	service.pm.AddWeb("api_web_socket_client", webServer)

	return service, nil
}

// Start ...
func (c *Client) Start(waitGroup ...*sync.WaitGroup) error {
	var wg *sync.WaitGroup

	if len(waitGroup) == 0 {
		wg = &sync.WaitGroup{}
		wg.Add(1)
	} else {
		wg = waitGroup[0]
	}

	defer wg.Done()

	err := c.pm.Start()

	if err == nil {
		c.started = true
	}

	return err
}

// Stop ...
func (c *Client) Stop(waitGroup ...*sync.WaitGroup) error {
	var wg *sync.WaitGroup

	if len(waitGroup) == 0 {
		wg = &sync.WaitGroup{}
		wg.Add(1)
	} else {
		wg = waitGroup[0]
	}

	defer wg.Done()

	err := c.pm.Stop()

	if err == nil {
		c.started = false
	}

	return err
}

func (c *Client) Started() bool {
	return c.started
}

// Subscribe ...
func (c *Client) Subscribe(topic, channel string) error {
	request, err := c.client.NewRequest(web.MethodPut, fmt.Sprintf("%s/subscribe/%s/%s", c.config.ServerAddress, topic, channel))
	if err != nil {
		return err
	}

	request.SetHeader(HeaderGatewayKey, []string{c.server.GetAddress()})

	response, err := request.Send()
	if err != nil {
		return err
	}

	fmt.Printf("\nserver response %s", string(response.Body))

	return nil
}

// Unsubscribe ...
func (c *Client) Unsubscribe(topic, channel string) error {
	request, err := c.client.NewRequest(web.MethodDelete, fmt.Sprintf("%s/unsubscribe/%s/%s", c.config.ServerAddress, topic, channel))
	if err != nil {
		return err
	}

	request.SetHeader(HeaderGatewayKey, []string{c.server.GetAddress()})

	response, err := request.Send()
	if err != nil {
		return err
	}

	fmt.Printf("\nserver response %s", string(response.Body))

	return nil
}

// Publish ...
func (c *Client) Publish(topic, channel string, message []byte) error {
	request, err := c.client.NewRequest(web.MethodPost, fmt.Sprintf("%s/new-message/%s/%s", c.config.ServerAddress, topic, channel))
	if err != nil {
		return err
	}

	request.SetHeader(HeaderGatewayKey, []string{c.server.GetAddress()})

	response, err := request.WithBody(message, web.ContentTypeApplicationJSON).Send()
	if err != nil {
		return err
	}

	fmt.Printf("\nserver response %s", string(response.Body))

	return nil
}

// Listen ...
func (c *Client) Listen(topic, channel string, handler MessageHandler) {
	mapChannels, ok := c.listeners[topic]
	if !ok {
		mapChannels = make(map[string]MessageHandler)
		c.listeners[topic] = mapChannels
	}

	mapChannels[channel] = handler
}

// Forget ...
func (c *Client) Forget(topic, channel string, handler MessageHandler) {
	if mapChannels, ok := c.listeners[topic]; ok {
		delete(mapChannels, channel)
	}
}

// Forget ...
func (c *Client) Wait() {
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)

	select {
	case <-termChan:
		c.Stop()
		c.logger.Infof("received term signal")
	}
}

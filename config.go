package socket

import (
	"fmt"
	"github.com/joaosoft/manager"
)

// AppConfig ...
type AppConfig struct {
	Socket SocketConfig `json:"socket"`
}

// SocketConfig ...
type SocketConfig struct {
	Server *ServerConfig `json:"server"`
	Client *ClientConfig `json:"client"`
}

// NewConfig ...
func NewConfig() (*AppConfig, manager.IConfig, error) {
	appConfig := &AppConfig{}
	simpleConfig, err := manager.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", GetEnv()), appConfig)

	return appConfig, simpleConfig, err
}

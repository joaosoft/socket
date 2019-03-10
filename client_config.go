package socket

import "github.com/joaosoft/web"

// ClientConfig ...
type ClientConfig struct {
	ServerAddress string `json:"server_address"`
	*web.ClientConfig
}

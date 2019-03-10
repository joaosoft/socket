package socket

import (
	"github.com/joaosoft/web"
)

func (c *clientController) registerRoutes() error {
	return c.server.AddRoutes(
		web.NewRoute(web.MethodPost, "/api/v1/new-message/:topic/:channel", c.handleNewMessage),
	)
}

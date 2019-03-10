package socket

import (
	"github.com/joaosoft/web"
)

func (c *serverController) registerRoutes() error {
	return c.server.AddRoutes(
		web.NewRoute(web.MethodPost, "/api/v1/new-message/:topic/:channel", c.handleNewMessage),
		web.NewRoute(web.MethodPut, "/api/v1/subscribe/:topic/:channel", c.handleSubscribe),
		web.NewRoute(web.MethodDelete, "/api/v1/unsubscribe/:topic/:channel", c.handleUnsubscribe),
	)
}

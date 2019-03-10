package socket

import (
	"github.com/joaosoft/web"
)

type clientController struct {
	*Client
}

func (c *Client) initController() *clientController {
	clientController := &clientController{Client: c}
	clientController.registerRoutes()

	return clientController
}

func (c *clientController) handleNewMessage(ctx *web.Context) error {
	topic := ctx.Request.GetUrlParam("topic")
	channel := ctx.Request.GetUrlParam("channel")

	c.logger.Infof("received a new message from server [topic: %s, channel: %s, body: %s]",
		topic,
		channel,
		string(ctx.Request.Body))

	if mapTopic, ok := c.listeners[topic]; ok {
		if mapChannel, ok := mapTopic[channel]; ok {
			c.logger.Infof("client handled the message [topic: %s, channel: %s, body: %s]",
				topic,
				channel,
				string(ctx.Request.Body))
			ctx.Response.JSON(web.StatusOK, &messageAcknowledge{Acknowledge: true, Errors: []error{mapChannel(ctx.Request.Body)}})
		}
	}

	return ctx.Response.JSON(web.StatusOK, &messageAcknowledge{Acknowledge: false})
}

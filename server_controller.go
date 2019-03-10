package socket

import (
	"fmt"

	"github.com/joaosoft/web"
)

type serverController struct {
	*Server
}

func (s *Server) initController() *serverController {
	serverController := &serverController{Server: s}
	serverController.registerRoutes()

	return serverController
}

func (s *serverController) handleNewMessage(ctx *web.Context) error {
	topic := ctx.Request.GetUrlParam("topic")
	channel := ctx.Request.GetUrlParam("channel")

	s.logger.Infof("received a new message from client [topic: %s, channel: %s, body: %s]",
		topic,
		channel,
		string(ctx.Request.Body))

	if mapTopic, ok := s.listeners[topic]; ok {
		if mapChannel, ok := mapTopic[channel]; ok {
			for _, listener := range mapChannel {
				request, err := s.client.NewRequest(web.MethodPost, fmt.Sprintf("%s/api/v1/new-message/%s/%s", listener.gateway, topic, channel))
				if err != nil {
					return ctx.Response.JSON(web.StatusBadRequest, &messageAcknowledge{
						Acknowledge: true,
						Error:       err,
					})
				}

				response, err := request.WithBody(ctx.Request.Body, web.ContentTypeApplicationJSON).Send()
				if err != nil {
					return ctx.Response.JSON(web.StatusBadRequest, &messageAcknowledge{
						Acknowledge: true,
						Error:       err,
					})
				}

				fmt.Printf("\nserver sending message to client [topic: %s, channel: %s message: %s, gateway: %s]", topic, channel, string(response.Body), listener.gateway)

				return ctx.Response.JSON(web.StatusOK, &messageAcknowledge{
					Acknowledge: true,
				})
			}
		}
	}

	return ctx.Response.JSON(web.StatusOK, &messageAcknowledge{
		Acknowledge: false,
	})
}

func (s *serverController) handleSubscribe(ctx *web.Context) error {
	topic := ctx.Request.GetUrlParam("topic")
	channel := ctx.Request.GetUrlParam("channel")
	gateway := ctx.Request.GetHeader("gateway")

	s.logger.Infof("subscribe new client [topic: %s, channel: %s, gateway: %s]",
		topic,
		channel,
		gateway)

	mapChannels, ok := s.listeners[topic]
	if !ok {
		mapChannels = make(map[string]map[string]*listener)
		s.listeners[topic] = mapChannels
	}

	mapListeners, ok := mapChannels[channel]
	if !ok {
		mapListeners = make(map[string]*listener)
		mapChannels[channel] = mapListeners
	}

	mapListeners[gateway] = &listener{gateway: gateway}
	mapListeners[gateway] = &listener{gateway: gateway}

	return ctx.Response.JSON(web.StatusOK, &messageAcknowledge{
		Acknowledge: true,
	})
}

func (s *serverController) handleUnsubscribe(ctx *web.Context) error {
	topic := ctx.Request.GetUrlParam("topic")
	channel := ctx.Request.GetUrlParam("channel")
	gateway := ctx.Request.GetHeader("gateway")

	s.logger.Infof("unsubscribe new client [topic: %s, channel: %s, gateway: %s]",
		topic,
		channel,
		gateway)

	if mapChannels, ok := s.listeners[topic]; ok {
		if mapListeners, ok := mapChannels[channel]; ok {
			delete(mapListeners, gateway)
		}
	}

	return ctx.Response.JSON(web.StatusOK, &messageAcknowledge{
		Acknowledge: true,
	})
}

package api

import (
	"abrarvan_challenge/api/middleware"
	customRouter "abrarvan_challenge/api/router"
	"abrarvan_challenge/config"
	"abrarvan_challenge/infrastructure/persistance/broker"
	"abrarvan_challenge/logging"

	"github.com/gin-gonic/gin"
)

type WebService struct {
	cfg     *config.Config
	logger  logging.Logger
	handler *MessageHandler
}

type MessageHandler struct {
	logger logging.Logger
}

func NewWebService(cfg *config.Config, logger logging.Logger) *WebService {
	handler := &MessageHandler{logger: logger}
	return &WebService{cfg: cfg, logger: logger, handler: handler}
}

func (ws *WebService) SetupRouter() *gin.Engine {
	gin.SetMode(config.GetConfig().Server.RunMode)
	router := gin.New()
	router.Use(middleware.DefaultStructuredLogger(ws.cfg))
	_, err := broker.CreateChannel("producerChannel", "my_queue", broker.WithDurable(true), broker.WithAutoDelete(false))
	if err != nil {
		ws.logger.Fatal(logging.RabbitMQ, logging.Startup, "Failed to create producer channel: "+err.Error(), nil)
	}
	return router
}

func (ws *WebService) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")
	healthCheckApiGroup := api.Group("/health")
    sendSmsApiGroup := api.Group("/notificaiton")
	customRouter.Health(healthCheckApiGroup)
    customRouter.SendSmsRouter(sendSmsApiGroup)
}

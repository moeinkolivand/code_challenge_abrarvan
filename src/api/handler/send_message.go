package handler

import (
	"abrarvan_challenge/config"
	"abrarvan_challenge/infrastructure/persistance/broker"
	"abrarvan_challenge/logging"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

var logger logging.Logger = logging.NewLogger(config.GetConfig())

type SendMessageHandler struct {
}

func NewSendMessageHandler() *SendMessageHandler {
	return &SendMessageHandler{}
}

func (sendMessageHandler *SendMessageHandler) SendMessageHandler(c *gin.Context) {
	var request struct {
		Message string `json:"message"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Error(logging.WebService, logging.Request, "Invalid request payload: "+err.Error(), nil)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	messageBody, err := json.Marshal(map[string]string{"message": request.Message})
	if err != nil {
		logger.Error(logging.RabbitMQ, logging.Publish, "Failed to marshal message: "+err.Error(), nil)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process message"})
		return
	}

	err = broker.Publish("producerChannel", "", "my_queue", messageBody)
	if err != nil {
		logger.Error(logging.RabbitMQ, logging.Publish, "Failed to publish message: "+err.Error(), nil)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish message"})
		return
	}

	logger.Info(logging.RabbitMQ, logging.Publish, "Message published successfully", map[logging.ExtraKey]interface{}{
		"body": string(messageBody),
	})
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Message queued for processing",
	})
}

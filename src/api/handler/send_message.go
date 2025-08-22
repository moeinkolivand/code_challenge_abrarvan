package handler

import (
	dto "abrarvan_challenge/api/dto"
	"abrarvan_challenge/config"
	"abrarvan_challenge/logging"
	"abrarvan_challenge/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

var logger logging.Logger = logging.NewLogger(config.GetConfig())

type UserApiHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserApiHandler {
	return &UserApiHandler{
		userService: userService,
	}
}

func (userHandler *UserApiHandler) UserHandler(c *gin.Context) {
	var request dto.SendMessageRequestDto

	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Error(logging.WebService, logging.Request, "Invalid request payload: "+err.Error(), nil)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err := userHandler.userService.SendMessage(request.PhoneNumber, request.Message)
	if err != nil {
		logger.Error(logging.WebService, logging.Request, err.Error(), nil)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized user"})
		return
	}

	logger.Info(logging.RabbitMQ, logging.Publish, "Message published successfully", map[logging.ExtraKey]interface{}{
		"body": "Hello From Api",
	})

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Message queued for processing",
	})
}

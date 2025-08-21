package router

import (
	"abrarvan_challenge/api/handler"

	"github.com/gin-gonic/gin"
)

func SendSmsRouter(r *gin.RouterGroup) {
	handler := handler.NewSendMessageHandler()
	r.POST("/send-message", handler.SendMessageHandler)
}

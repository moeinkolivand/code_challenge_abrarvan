package router

import (
	"abrarvan_challenge/api/handler"
	"abrarvan_challenge/infrastructure/persistance/repository"
	"abrarvan_challenge/provider"
	usrService "abrarvan_challenge/service"
	"gorm.io/gorm"
	"math/rand"

	"github.com/gin-gonic/gin"
)

func SendSmsRouter(r *gin.RouterGroup, db *gorm.DB) {
	userRepository := repository.NewUserRepository(db)
	providerName := []string{"provider_one", "provider_two"}[rand.Intn(2)]
	smsProvider := provider.NewProvider(providerName)
	userService := usrService.NewUserService(smsProvider, userRepository)
	userHandler := handler.NewUserHandler(userService)
	r.POST("/send-message", userHandler.UserHandler)
}

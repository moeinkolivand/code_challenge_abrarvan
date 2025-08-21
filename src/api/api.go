package api

import (
	"abrarvan_challenge/config"
	"abrarvan_challenge/logging"

	"github.com/gin-gonic/gin"
)


var logger = logging.NewLogger(config.GetConfig())


func InitServer(cfg *config.Config) *gin.Engine {
	webService := NewWebService(cfg, logger)
	router := webService.SetupRouter()
	webService.RegisterRoutes(router)
	return router
}
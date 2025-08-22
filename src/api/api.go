package api

import (
	"abrarvan_challenge/config"
	"abrarvan_challenge/logging"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

var logger = logging.NewLogger(config.GetConfig())

func InitServer(cfg *config.Config, db *gorm.DB) *gin.Engine {
	webService := NewWebService(cfg, logger, db)
	router := webService.SetupRouter()
	webService.RegisterRoutes(router)
	return router
}

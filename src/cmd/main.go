package main

import (
	"abrarvan_challenge/config"
	"abrarvan_challenge/infrastructure/cache"
	"abrarvan_challenge/infrastructure/persistance/broker"
	"abrarvan_challenge/infrastructure/persistance/database"
	"abrarvan_challenge/logging"
	"abrarvan_challenge/model"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	err := os.Setenv("APP_ENV", "local")
	if err != nil {
		return
	}
	cfg := config.GetConfig()
	logger := logging.NewLogger(cfg)

	err = cache.InitRedis(cfg)
	defer cache.CloseRedis()
	if err != nil {
		logger.Fatal(logging.Redis, logging.Startup, err.Error(), nil)
	}

	err = database.InitDb(cfg)
	defer database.CloseDb()
	if err != nil {
		logger.Fatal(logging.Postgres, logging.Startup, err.Error(), nil)
	}
	//migration.Up1()
	err = broker.InitRabbitMq(cfg)
	if err != nil {
		logger.Fatal(logging.RabbitMQ, logging.Startup, err.Error(), nil)
	}
	err = model.MigrateDatabaseTables(database.GetDb())
	if err != nil {
		logger.Fatal(logging.DatabaseMigration, logging.Startup, err.Error(), nil)
	}
	//api.InitServer(cfg)
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	runRrr := router.Run()
	if runRrr != nil {
		return
	} // listen and serve on 0.0.0.0:8080
}

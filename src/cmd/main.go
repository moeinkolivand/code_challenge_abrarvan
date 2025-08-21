package main

import (
	"abrarvan_challenge/config"
	"abrarvan_challenge/infrastructure/cache"
	"abrarvan_challenge/infrastructure/persistance/broker"
	"abrarvan_challenge/infrastructure/persistance/database"
	"abrarvan_challenge/logging"
	"abrarvan_challenge/model"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	err := os.Setenv("APP_ENV", "local")
	if err != nil {
		return
	}
	cfg := config.GetConfig()
	logger := logging.NewLogger(cfg)
	runConsumer := flag.Bool("consumer", false, "Run the message consumer")

	// Queue options
	queueName := flag.String("queue", "", "Queue name to consume from")
	queueDurable := flag.Bool("queue-durable", true, "Queue durable")
	queueAutoDelete := flag.Bool("queue-autodelete", false, "Queue auto-delete")
	queueExclusive := flag.Bool("queue-exclusive", false, "Queue exclusive")
	queueNoWait := flag.Bool("queue-nowait", false, "Queue nowait")

	// Consumer options
	consumerName := flag.String("consumer-name", "", "Consumer name")
	autoAck := flag.Bool("auto-ack", true, "Auto acknowledge messages")
	consExclusive := flag.Bool("cons-exclusive", false, "Exclusive consumer")
	noLocal := flag.Bool("no-local", false, "No local messages")
	consNoWait := flag.Bool("cons-nowait", false, "Consumer noWait")

	flag.Parse()

	if *runConsumer {
		runConsumerMode(cfg, logger, *queueName, broker.ConsumeOptions{
			Consumer:  *consumerName,
			AutoAck:   *autoAck,
			Exclusive: *consExclusive,
			NoLocal:   *noLocal,
			NoWait:    *consNoWait,
			Args:      nil,
		}, []broker.QueueOption{
			broker.WithDurable(*queueDurable),
			broker.WithAutoDelete(*queueAutoDelete),
			broker.WithExclusive(*queueExclusive),
			broker.WithNoWait(*queueNoWait),
		})
	} else {
		runWebServiceMode(cfg)
	}

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

}
func runConsumerMode(cfg *config.Config, log logging.Logger, queueName string, consOpts broker.ConsumeOptions, queueOpts []broker.QueueOption) {
	fmt.Println(1)
	log.Info(logging.RabbitMQ, logging.Startup, "Starting consumer mode...", nil)

	// Initialize RabbitMQ connection
	if err := broker.InitRabbitMq(cfg); err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	_, err := broker.CreateChannel("consumerChannel", queueName, queueOpts...)
	if err != nil {
		log.Fatalf("Failed to create channel: %v", err)
	}

	msgs, err := broker.Consume("consumerChannel", queueName, consOpts)
	if err != nil {
		log.Fatalf("Failed to start consuming: %v", err)
	}
	fmt.Println(123456)
	// Handle OS signals for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	log.Info(logging.RabbitMQ, logging.Startup, "Consumer started. Waiting for messages...", nil)

	for {
		select {
		case msg := <-msgs:
			log.Infof("Received message from custom log: %s", msg.Body)
			// Add your processing logic here
		case <-sigChan:
			log.Infof("Shutting down consumer...")
			return
		}
	}
}

func runWebServiceMode(cfg *config.Config) {
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

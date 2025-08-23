package runner

import (
	"abrarvan_challenge/api"
	"abrarvan_challenge/config"
	"abrarvan_challenge/infrastructure/cache"
	"abrarvan_challenge/infrastructure/persistance/broker"
	"abrarvan_challenge/infrastructure/persistance/database"
	"abrarvan_challenge/logging"
	"abrarvan_challenge/model"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"gorm.io/gorm"
)

type App struct {
	cfg    *config.Config
	Logger logging.Logger
	db     *gorm.DB
}

func NewApp() (*App, error) {
	// Load configuration
	cfg := config.GetConfig()
	logger := logging.NewLogger(cfg)

	// Initialize Redis
	if err := cache.InitRedis(cfg); err != nil {
		logger.Fatal(logging.Redis, logging.Startup, "Failed to initialize Redis: "+err.Error(), nil)
		return nil, err
	}

	// Initialize PostgreSQL
	if err := database.InitDb(cfg); err != nil {
		logger.Fatal(logging.Postgres, logging.Startup, "Failed to initialize database: "+err.Error(), nil)
		return nil, err
	}
	databaseConnection := database.GetDb()
	err := model.SeedUsers(databaseConnection)
	if err != nil {
		return nil, err
	}

	// Initialize RabbitMQ
	if err := broker.InitRabbitMq(cfg); err != nil {
		logger.Fatal(logging.RabbitMQ, logging.Startup, "Failed to initialize RabbitMQ: "+err.Error(), nil)
		return nil, err
	}

	// Migrate database tables
	if err := model.MigrateDatabaseTables(databaseConnection); err != nil {
		logger.Fatal(logging.DatabaseMigration, logging.Startup, "Failed to migrate database: "+err.Error(), nil)
		return nil, err
	}

	return &App{cfg: cfg, Logger: logger, db: databaseConnection}, nil
}

func ParseFlags() (bool, string, broker.ConsumeOptions, []broker.QueueOption) {
	runConsumer := flag.Bool("consumer", false, "Run the message consumer")
	queueName := flag.String("queue", "", "Queue name to consume from")
	queueDurable := flag.Bool("queue-durable", true, "Queue durable")
	queueAutoDelete := flag.Bool("queue-autodelete", false, "Queue auto-delete")
	queueExclusive := flag.Bool("queue-exclusive", false, "Queue exclusive")
	queueNoWait := flag.Bool("queue-nowait", false, "Queue no-wait")
	consumerName := flag.String("consumer-name", "", "Consumer name")
	autoAck := flag.Bool("auto-ack", true, "Auto acknowledge messages")
	consExclusive := flag.Bool("cons-exclusive", false, "Exclusive consumer")
	noLocal := flag.Bool("no-local", false, "No local messages")
	consNoWait := flag.Bool("cons-nowait", false, "Consumer no-wait")

	flag.Parse()

	consOpts := broker.ConsumeOptions{
		Consumer:  *consumerName,
		AutoAck:   *autoAck,
		Exclusive: *consExclusive,
		NoLocal:   *noLocal,
		NoWait:    *consNoWait,
		Args:      nil,
	}

	queueOpts := []broker.QueueOption{
		broker.WithDurable(*queueDurable),
		broker.WithAutoDelete(*queueAutoDelete),
		broker.WithExclusive(*queueExclusive),
		broker.WithNoWait(*queueNoWait),
	}

	return *runConsumer, *queueName, consOpts, queueOpts
}

func RunConsumerMode(app *App, queueName string, consOpts broker.ConsumeOptions, queueOpts []broker.QueueOption) {
	app.Logger.Info(logging.RabbitMQ, logging.Startup, "Starting consumer mode...", nil)

	ch, err := broker.CreateChannel("consumerChannel", queueName, queueOpts...)
	if err != nil {
		app.Logger.Fatal(logging.RabbitMQ, logging.Startup, "Failed to create channel: "+err.Error(), nil)
	}

	msgs, err := broker.Consume("consumerChannel", queueName, consOpts)
	if err != nil {
		app.Logger.Fatal(logging.RabbitMQ, logging.Startup, "Failed to start consuming: "+err.Error(), nil)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	app.Logger.Info(logging.RabbitMQ, logging.Startup, "Consumer started. Waiting for messages...", nil)

	for {
		select {
		case msg := <-msgs:
			app.Logger.Info(logging.RabbitMQ, logging.MessageReceived, "Received message from custom log", map[logging.ExtraKey]interface{}{
				"body": string(msg.Body),
			})
			if !consOpts.AutoAck {
				if err := msg.Ack(false); err != nil {
					app.Logger.Error(logging.RabbitMQ, logging.MessageAck, "Failed to acknowledge message: "+err.Error(), nil)
				} else {
					app.Logger.Info(logging.RabbitMQ, logging.MessageAck, "Message acknowledged", nil)
				}
			}
		case <-sigChan:
			app.Logger.Info(logging.RabbitMQ, logging.Shutdown, "Shutting down consumer...", nil)
			if ch != nil {
				ch.Close()
			}
			cache.CloseRedis()
			database.CloseDb()
			return
		}
	}
}

func RunWebServiceMode(app *App) {
	router := api.InitServer(app.cfg, app.db)
	addr := ":" + app.cfg.Server.InternalPort
	if err := router.Run(addr); err != nil {
		app.Logger.Fatal(logging.WebService, logging.Startup, "Failed to start web server: "+err.Error(), nil)
	}
}

package main

import (
	"abrarvan_challenge/cmd/runner"
	"abrarvan_challenge/logging"
	"fmt"
	"os"
)

func main() {
	app, err := runner.NewApp()
	if err != nil {
		fmt.Printf("Failed to initialize app: %v\n", err)
		os.Exit(1)
	}

	runConsumer, queueName, consOpts, queueOpts := runner.ParseFlags()
	if runConsumer {
		if queueName == "" {
			app.Logger.Fatal(logging.RabbitMQ, logging.Startup, "Queue name is required in consumer mode", nil)
		}
		runner.RunConsumerMode(app, queueName, consOpts, queueOpts)
	} else {
		runner.RunWebServiceMode(app)
	}
}

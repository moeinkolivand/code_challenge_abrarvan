package broker

import (
	"abrarvan_challenge/config"
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var rabbitConn *amqp.Connection
var channels = make(map[string]*amqp.Channel)

// QueueOption defines a functional option for queue declaration
type QueueOption func(*queueOptions)

type queueOptions struct {
	durable    bool
	autoDelete bool
	exclusive  bool
	noWait     bool
}

// Consumer represents a RabbitMQ consumer.
type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
}

type ConsumeOptions struct {
	Consumer  string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqp.Table
}

// defaultConsumeOptions provides default values
var defaultConsumeOptions = ConsumeOptions{
	Consumer:  "",
	AutoAck:   true,
	Exclusive: false,
	NoLocal:   false,
	NoWait:    false,
	Args:      nil,
}

func WithDurable(d bool) QueueOption {
	return func(o *queueOptions) { o.durable = d }
}
func WithAutoDelete(ad bool) QueueOption {
	return func(o *queueOptions) { o.autoDelete = ad }
}
func WithExclusive(e bool) QueueOption {
	return func(o *queueOptions) { o.exclusive = e }
}
func WithNoWait(nw bool) QueueOption {
	return func(o *queueOptions) { o.noWait = nw }
}

func InitRabbitMq(cfg *config.Config) error {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.RabbitMQ.User,
		cfg.RabbitMQ.Password,
		cfg.RabbitMQ.Host,
		cfg.RabbitMQ.Port,
	)
	connection, err := amqp.Dial(url)
	if err != nil {
		return err
	}
	rabbitConn = connection
	return nil
}

func GetRabbitConnection() *amqp.Connection {
	return rabbitConn
}

func CreateChannel(name, queueName string, opts ...QueueOption) (*amqp.Channel, error) {
	if rabbitConn == nil {
		return nil, fmt.Errorf("rabbitMQ connection is nil")
	}

	ch, err := rabbitConn.Channel()
	if err != nil {
		fmt.Println("Failed to create channel: %v", err)
		return nil, err
	}

	// Default options
	qOpts := &queueOptions{
		durable:    true,
		autoDelete: false,
		exclusive:  false,
		noWait:     false,
	}

	// Apply functional options
	for _, opt := range opts {
		opt(qOpts)
	}

	// Declare the queue
	_, err = ch.QueueDeclare(
		queueName,
		qOpts.durable,
		qOpts.autoDelete,
		qOpts.exclusive,
		qOpts.noWait,
		nil,
	)
	if err != nil {
		return nil, err
	}
	SetChannel(name, ch)
	return ch, nil
}

func GetChannel(name string) (*amqp.Channel, error) {
	ch, ok := channels[name]
	if !ok {
		return nil, fmt.Errorf("channel %s not found", name)
	}
	return ch, nil
}

func Publish(channelKey, exchange, routingKey string, body []byte) error {
	ch, err := GetChannel(channelKey)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return ch.PublishWithContext(ctx, exchange, routingKey, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}

func Consume(channelName, queueName string, opts ...ConsumeOptions) (<-chan amqp.Delivery, error) {
	ch, err := GetChannel(channelName)
	if err != nil {
		fmt.Println(123123123)
		return nil, err
	}

	if rabbitConn == nil {
		return nil, amqp.ErrClosed
	}
	finalOpts := defaultConsumeOptions

	// If user passed options, overwrite defaults
	if len(opts) > 0 {
		opt := opts[0]
		if opt.Consumer != "" {
			finalOpts.Consumer = opt.Consumer
		}
		finalOpts.AutoAck = opt.AutoAck
		finalOpts.Exclusive = opt.Exclusive
		finalOpts.NoLocal = opt.NoLocal
		finalOpts.NoWait = opt.NoWait
		if opt.Args != nil {
			finalOpts.Args = opt.Args
		}
	}

	return ch.Consume(
		queueName,
		finalOpts.Consumer,
		finalOpts.AutoAck,
		finalOpts.Exclusive,
		finalOpts.NoLocal,
		finalOpts.NoWait,
		finalOpts.Args,
	)
}

func (c *Consumer) Close() error {
	if err := c.channel.Close(); err != nil {
		return err
	}
	return c.conn.Close()
}

// SetChannel stores a channel with a given name
func SetChannel(name string, ch *amqp.Channel) {
	channels[name] = ch
}

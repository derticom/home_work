package ampq

import (
	"fmt"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	exchange    = "test-exchange"   // Durable, non-auto-deleted AMQP exchange name
	queueName   = "test-queue"      // Ephemeral AMQP queue name
	bindingKey  = "test-key"        // AMQP binding key
	consumerTag = "simple-consumer" // AMQP consumer tag (should not be blank)
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	done    chan error
	log     *slog.Logger
}

// NewConsumer создает новый экземпляр Consumer.
func NewConsumer(log *slog.Logger) *Consumer {
	return &Consumer{
		tag:  consumerTag,
		done: make(chan error),
		log:  log,
	}
}

// Setup создает соединение, канал, объявляет обменник, очередь и привязывает очередь к обменнику.
func (c *Consumer) Setup(amqpURI string) error {
	var err error

	c.log.Debug("dialing", "amqpURI", amqpURI)
	c.conn, err = amqp.Dial(amqpURI)
	if err != nil {
		return fmt.Errorf("failed to amqp.Dial: %w", err)
	}

	go func() {
		err := <-c.conn.NotifyClose(make(chan *amqp.Error))
		c.log.Error("connection closed", "error", err)
	}()

	c.log.Debug("got Connection, getting Channel")
	c.channel, err = c.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to conn.Channel: %w", err)
	}

	c.log.Debug("got Channel, declaring Exchange", "exchange", exchange)
	if err = c.channel.ExchangeDeclare(
		exchange, // name of the exchange
		"direct", // type (example: "direct", "fanout", "topic")
		true,     // durable
		false,    // delete when complete
		false,    // internal
		false,    // noWait
		nil,      // arguments
	); err != nil {
		return fmt.Errorf("failed to channel.ExchangeDeclare: %w", err)
	}

	c.log.Debug("declared Exchange, declaring Queue", "queue name", queueName)
	queue, err := c.channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to channel.QueueDeclare: %w", err)
	}

	c.log.Debug("declared Queue ", "name", queue.Name)

	if err = c.channel.QueueBind(
		queue.Name, // name of the queue
		bindingKey, // bindingKey
		exchange,   // sourceExchange
		false,      // noWait
		nil,        // arguments
	); err != nil {
		return fmt.Errorf("failed to channel.QueueBind: %w", err)
	}

	return nil
}

// Consume запускает прослушивание очереди и возвращает канал для получения сообщений.
func (c *Consumer) Consume() (<-chan amqp.Delivery, error) {
	c.log.Debug("Queue bound to Exchange, starting Consume", "consumer tag", c.tag)
	deliveries, err := c.channel.Consume(
		queueName, // name
		c.tag,     // consumerTag,
		false,     // noAck
		false,     // exclusive
		false,     // noLocal
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to channel.Consume: %w", err)
	}

	return deliveries, nil
}

// Shutdown закрывает соединение и канал.
func (c *Consumer) Shutdown() error {
	if err := c.channel.Cancel(c.tag, true); err != nil {
		return fmt.Errorf("failed to channel.Cancel: %w", err)
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("failed to conn.Close: %w", err)
	}

	c.log.Info("AMQP shutdown OK")
	return nil
}

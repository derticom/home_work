package ampq

import (
	"fmt"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	exchangeName = "test-exchange" // Durable AMQP exchange name
	exchangeType = "direct"        // Exchange type - direct|fanout|topic|x-custom
	routingKey   = "test-key"      // AMQP routing key
)

type Producer struct {
	log *slog.Logger

	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewProducer(uri string, log *slog.Logger) (*Producer, error) {
	connection, err := amqp.Dial(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to amqp.Dial: %w", err)
	}

	channel, err := connection.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to connection.Channel: %w", err)
	}

	log.Debug("got Channel, declaring", "Exchange", exchangeName)
	if err := channel.ExchangeDeclare(
		exchangeName, // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return nil, fmt.Errorf("failed to channel.ExchangeDeclare: %w", err)
	}

	return &Producer{
		connection: connection,
		channel:    channel,
		log:        log,
	}, nil
}

func (p *Producer) Publish(body []byte) error {
	p.log.Debug("declared Exchange, publishing", "body", body)
	if err := p.channel.Publish(
		exchangeName, // publish to an exchange
		routingKey,   // routing to 0 or more queues
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            body,
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		return fmt.Errorf("failed to channel.Publish: %w", err)
	}

	return nil
}

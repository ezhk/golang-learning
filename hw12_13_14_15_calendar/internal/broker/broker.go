package broker

import (
	"encoding/json"
	"fmt"

	"github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/config"
	"github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/structs"
	"github.com/streadway/amqp"
)

type Broker struct {
	Cfg *config.Configuration

	Conn    *amqp.Connection
	Channel *amqp.Channel
}

type Producer struct {
	Broker
}

type Consumer struct {
	Broker
	Queue *amqp.Queue
}

func (b *Broker) Connect() error {
	connection, err := amqp.Dial(b.Cfg.Broker.AMQP)
	if err != nil {
		return fmt.Errorf("AMQP dial error: %w", err)
	}
	b.Conn = connection

	return nil
}

func (b *Broker) Close() error {
	return b.Conn.Close()
}

func (b *Broker) ChannelDeclare() error {
	channel, err := b.Conn.Channel()
	if err != nil {
		return fmt.Errorf("channel declare error: %w", err)
	}
	b.Channel = channel

	return nil
}

func (b *Broker) ExchangeDeclare() error {
	return b.Channel.ExchangeDeclare(
		b.Cfg.Broker.Exchange,     // name
		b.Cfg.Broker.ExchangeType, // type
		true,                      // durable
		false,                     // auto-deleted
		false,                     // internal
		false,                     // noWait
		nil,                       // arguments
	)
}

func NewProducer(cfg *config.Configuration) *Producer {
	return &Producer{Broker: Broker{Cfg: cfg}}
}

func (p *Producer) Init() error {
	if err := p.Connect(); err != nil {
		return err
	}

	if err := p.ChannelDeclare(); err != nil {
		return err
	}
	if err := p.ExchangeDeclare(); err != nil {
		return err
	}

	return nil
}

func (p *Producer) Publish(e structs.Event) error {
	payload, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("marshall event error: %w", err)
	}

	return p.Channel.Publish(
		p.Cfg.Broker.Exchange,   // publish to an exchange
		p.Cfg.Broker.RoutingKey, // routing to 0 or more queues
		false,                   // mandatory
		false,                   // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "application/json",
			ContentEncoding: "UTF-8",
			Body:            payload,
			DeliveryMode:    amqp.Persistent,
		},
	)
}

func NewConsumer(cfg *config.Configuration) *Consumer {
	return &Consumer{Broker: Broker{Cfg: cfg}, Queue: nil}
}

// Init create connection and define rabbitMQ objects: exchange and queue.
func (c *Consumer) Init() error {
	if err := c.Connect(); err != nil {
		return err
	}

	// Producer and Consumer dec
	if err := c.ChannelDeclare(); err != nil {
		return err
	}
	if err := c.ExchangeDeclare(); err != nil {
		return err
	}

	// Consumer specific actions.
	if err := c.QueueDeclare(); err != nil {
		return err
	}
	if err := c.QueueBind(); err != nil {
		return err
	}

	return nil
}

func (c *Consumer) QueueDeclare() error {
	queueName := fmt.Sprintf("%s-queue", c.Cfg.Broker.RoutingKey)
	queue, err := c.Channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("queue declare error: %w", err)
	}
	c.Queue = &queue

	return nil
}

func (c *Consumer) QueueBind() error {
	return c.Channel.QueueBind(
		c.Queue.Name,            // name of the queue
		c.Cfg.Broker.RoutingKey, // bindingKey
		c.Cfg.Broker.Exchange,   // sourceExchange
		false,                   // noWait
		nil,                     // arguments
	)
}

func (c *Consumer) Consume() (<-chan amqp.Delivery, error) {
	msgCh, err := c.Channel.Consume(
		c.Queue.Name, // name
		"go-client",  // consumerTag,
		false,        // autoAck
		false,        // exclusive
		false,        // noLocal
		false,        // noWait
		nil,          // arguments
	)
	if err != nil {
		return msgCh, fmt.Errorf("consume error: %w", err)
	}

	return msgCh, nil
}

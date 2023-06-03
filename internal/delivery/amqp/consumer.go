package amqp

import (
	"context"

	"github.com/f0rmul/sensor-control/internal/models"
	"github.com/f0rmul/sensor-control/pkg/logger"
	snapshot_v1 "github.com/f0rmul/sensor-control/pkg/models"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"
)

type Service interface {
	PushAndSave(context.Context, *models.Snapshot) (*models.Snapshot, error)
}

type SnapshotConsumer struct {
	amqpConn *amqp.Connection
	service  Service
	logger   logger.Logger
}

func NewSnapshotConsumer(amqpConn *amqp.Connection, service Service, logger logger.Logger) *SnapshotConsumer {
	return &SnapshotConsumer{amqpConn: amqpConn, service: service, logger: logger}
}

func (c *SnapshotConsumer) CreateChannel(exchangeName, queueName, bindingKey, consumerTag string) (*amqp.Channel, error) {
	ch, err := c.amqpConn.Channel()

	if err != nil {
		return nil, errors.Wrap(err, "amqpConn.Channel()")
	}

	c.logger.Infof("Declaring exchange: %s", exchangeName)

	err = ch.ExchangeDeclare(
		exchangeName,
		exchangeKind,
		exchangeDurable, // persistence
		exchangeAutoDelete,
		exchangeInternal,
		exchangeNoWait,
		nil)

	if err != nil {
		return nil, errors.Wrap(err, "ch.Exchagedeclare()")
	}

	queue, err := ch.QueueDeclare(
		queueName,
		queueDurable,
		queueAutoDelete,
		queueExclusive,
		queueNoWait,
		nil)

	if err != nil {
		return nil, errors.Wrap(err, "ch.QueueDeclare()")
	}

	c.logger.Infof("Binding queue to exchange: Queue: %v , messageCount: %v ,"+
		"consumerCount: %v , exchange: %v , bindengKey: %v",
		queue.Name,
		queue.Messages,
		queue.Consumers,
		exchangeName,

		bindingKey)

	err = ch.QueueBind(
		queue.Name,
		bindingKey,
		exchangeName,
		queueNoWait,
		nil)

	if err != nil {
		return nil, errors.Wrap(err, "ch.QueueBind()")
	}

	err = ch.Qos(prefetchCount, prefetchSize, prefetchGlobal)

	if err != nil {
		return nil, errors.Wrap(err, "ch.Qos()")
	}

	return ch, nil
}

func (c *SnapshotConsumer) ProcessDelivery(ctx context.Context, deliveries <-chan amqp.Delivery) {

	for message := range deliveries {
		c.logger.Infof("Processing delivery: %v", message.DeliveryTag)

		p := new(snapshot_v1.Snapshot)

		if err := proto.Unmarshal(message.Body, p); err != nil {
			c.logger.Errorf("proto.Unmarshal(): %v", err)
			continue
		}

		snapshot := models.NewFromProto(p)

		saved, err := c.service.PushAndSave(ctx, snapshot)
		if err != nil {
			if err := message.Reject(false); err != nil {
				c.logger.Errorf("message.Reject(): %v", err)
			}
			c.logger.Error("failed to process message")
		} else {
			if err := message.Ack(false); err != nil {
				c.logger.Errorf("message.Ack(): %v", err)
			}
		}

		c.logger.Infof("Snapshot with ID: %s was consumed", saved.ID)
	}
	c.logger.Info("amqp.Delivery was closed")
}

func (c *SnapshotConsumer) StartConsumer(workerPool int, exchange, queueName, bindingKey, consumerTag string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, err := c.CreateChannel(exchange, queueName, bindingKey, consumerTag)

	if err != nil {
		return errors.Wrap(err, "c.CreateChannel()")
	}

	defer ch.Close()

	messages, err := ch.Consume(
		queueName,
		consumerTag,
		consumeAutoAck,
		consumeExclusive,
		consumeNoLocal,
		consumeNoWait,
		nil)

	if err != nil {
		return errors.Wrap(err, "ch.Consume()")
	}

	for i := 0; i < workerPool; i++ {
		go c.ProcessDelivery(ctx, messages)
	}

	chanErr := <-ch.NotifyClose(make(chan *amqp.Error))
	c.logger.Errorf("ch.NotifyClose(): %v", chanErr)
	return chanErr
}

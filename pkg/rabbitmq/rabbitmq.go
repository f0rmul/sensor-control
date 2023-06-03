package rabbitmq

import (
	"fmt"

	"github.com/f0rmul/sensor-control/config"
	"github.com/streadway/amqp"
)

func NewRabbitMQConn(cfg *config.Config) (*amqp.Connection, error) {
	connAddr := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		cfg.RabbitMQ.User,
		cfg.RabbitMQ.Password,
		cfg.RabbitMQ.Host,
		cfg.RabbitMQ.Port)
	return amqp.Dial(connAddr)
}

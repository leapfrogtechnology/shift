package mq

import (
	"fmt"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

// GetConnection returns the connection to rabbitmq
func GetConnection() *amqp.Connection {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")

	failOnError(err, "Failed to connect to RabbitMQ")

	return conn
}

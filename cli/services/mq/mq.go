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

func getQueue() (*amqp.Connection, *amqp.Channel, *amqp.Queue) {
	conn, err := amqp.Dial("amqp://shift:shiftdeveloper@dev-shiftmq.lftechnology.com:5672")

	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()

	failOnError(err, "Failed to open a channel")

	q, err := ch.QueueDeclare("Infrastructure",
		false, //durable bool
		false, //autoDelete bool
		false, //exclusive bool
		false, //noWait bool
		nil,   // args amqp.Table
	)

	failOnError(err, "Faied to declare a queue")

	return conn, ch, &q
}

// Publish message to queue
func Publish(message []byte) {
	conn, ch, q := getQueue()

	defer conn.Close()
	defer ch.Close()

	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        message,
	}

	ch.Publish("", q.Name, false, false, msg)
}

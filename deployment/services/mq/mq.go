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

	q, err := ch.QueueDeclare("Deployment",
		false, //durable bool
		false, //autoDelete bool
		false, //exclusive bool
		false, //noWait bool
		nil,   // args amqp.Table
	)

	failOnError(err, "Failed to declare a queue")

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

// Consume listens to message from queue
func Consume(deploy func([]byte)) {
	conn, ch, q := getQueue()

	defer conn.Close()
	defer ch.Close()
	msgs, err := ch.Consume(
		q.Name, // queue string
		"",     // consumer string
		true,   // autoAck bool
		false,  // exclusive bool
		false,  // noLocal bool
		false,  // noWait bool
		nil,    //args amqp.Table
	)

	failOnError(err, "Failed to register a consumer")

	for msg := range msgs {
		deploy(msg.Body)
	}
}

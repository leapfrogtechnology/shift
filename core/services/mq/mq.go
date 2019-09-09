package mq

import (
	"fmt"
	"github.com/leapfrogtechnology/shift/core/utils/logger"
	"github.com/streadway/amqp"
	"os"
	"time"
)

func GenerateQueueUrl() string {
	rabbitUserName := os.Getenv("RABBIT_USERNAME")
	rabbitPassword := os.Getenv("RABBIT_PASSWORD")
	rabbitHost := os.Getenv("RABBIT_HOST")
	rabbitPort := os.Getenv("RABBIT_PORT")
	rabbitUrl := fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitUserName, rabbitPassword, rabbitHost, rabbitPort)
	return rabbitUrl
}

func getQueue(ch *amqp.Channel, queueName string) *amqp.Queue {
	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	logger.FailOnError(err, "Failed to declare a queue: "+ queueName)
	return &q
}

func Consume(ch *amqp.Channel, qName string) <-chan amqp.Delivery {
	q := getQueue(ch, qName)
	messages, err := ch.Consume(
		q.Name, //name
		"",     //consumer
		false,  //autoAck
		false,  // exclusive
		false,  //noLocal
		false,  //noWait
		nil,    // args
	)
	logger.FailOnError(err, "Failed to Consume from Queue")
	return messages
}

func Publish(body []byte, ch *amqp.Channel, qName string) error {
	q := getQueue(ch, qName)
	err := ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			Headers:         nil,
			ContentType:     "application/json",
			ContentEncoding: "",
			DeliveryMode:    0,
			Priority:        0,
			CorrelationId:   "",
			ReplyTo:         "",
			Expiration:      "",
			MessageId:       "",
			Timestamp:       time.Time{},
			Type:            "",
			UserId:          "",
			AppId:           "",
			Body:            body,
		},
	)
	return err
}

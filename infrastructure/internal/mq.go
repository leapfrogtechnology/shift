package internal

import (
	"github.com/leapfrogtechnology/shift/infrastructure/utils"
	"github.com/streadway/amqp"
	"time"
)

func Publish(details string) error {
	conn, err := amqp.Dial(utils.GenerateQueueUrl())
	utils.FailOnError(err, "Failed to Connect to Message Queue")
	defer conn.Close()

	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"Deployment", // name
		false,        // durable
		false,        // delete when usused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	utils.FailOnError(err, "Failed to declare a queue")
	err = ch.Publish(
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
			Body:            []byte(details),
		},
	)
	if err != nil {
		utils.LogError(err, "Could not Publish output")
	}
	return err
}

package main

import (
	"github.com/leapfrogtechnology/shift/infrastructure/infrastructure"
	"github.com/leapfrogtechnology/shift/infrastructure/utils"
	"github.com/streadway/amqp"
	"log"
)

func main() {

	conn, err := amqp.Dial(utils.GenerateQueueUrl())
	utils.FailOnError(err, "Failed to Connect to Message Queue")
	defer conn.Close()

	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"initialize", // name
		false,   // durable
		false,   // delete when usused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	utils.FailOnError(err, "Failed to declare a queue")

	messages, err := ch.Consume(
		q.Name, //name
		"", //consumer
		true, //autoAck
		false, // exclusive
		false, //noLocal
		false, //noWait
		nil, // args
	)
	forever := make(chan bool)
	go func() {
		for message := range messages {
			log.Printf("Received a message: %s", message.MessageId)
			infrastructureInfo, err := infrastrucuture.InitializeFrontend(message.Body)
			if err != nil {
				utils.LogError(err, "Cannot Init Infrastructure")
			} else {
				utils.LogOutput(infrastructureInfo)
			}
		}
	}()
	utils.LogInfo("[*] Waiting for messages. To exit Press Ctrl+C")
	<-forever
}

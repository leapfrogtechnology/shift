package internals

import (
	"github.com/leapfrogtechnology/shift/core/services/mq"
	"github.com/leapfrogtechnology/shift/core/utils/logger"
	"github.com/leapfrogtechnology/shift/infrastructure/infrastructure"
	"github.com/streadway/amqp"
	"log"
)

func Initialize() {
	conn, err := amqp.Dial(mq.GenerateQueueUrl())
	logger.FailOnError(err, "Failed to Connect to Message Queue")
	defer conn.Close()
	ch, err := conn.Channel()
	logger.FailOnError(err, "Failed to open a channel")
	defer ch.Close()
	messages := mq.Consume(ch, "Infrastructure")
	forever := make(chan bool)
	go func() {
		for message := range messages {
			log.Printf("Received a message: %s", message.MessageId)
			infrastructureInfo, err := infrastrucuture.Initialize(message.Body)
			if err != nil {
				logger.LogError(err, "Cannot Init Infrastructure")
			} else {
				logger.LogOutput(infrastructureInfo)
				err = mq.Publish([]byte(infrastructureInfo), ch, "Deployment")
				if err != nil {
					logger.LogError(err, "Cannot Publish Output")
				}
			}
		}
	}()
	logger.LogInfo("[*] Waiting for messages. To exit Press Ctrl+C")
	<-forever
}

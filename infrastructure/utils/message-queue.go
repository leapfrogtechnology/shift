package utils

import (
	"fmt"
	"os"
)

func GenerateQueueUrl() string{
	rabbitUserName := os.Getenv("RABBIT_USERNAME")
	rabbitPassword := os.Getenv("RABBIT_PASSWORD")
	rabbitHost := os.Getenv("RABBIT_HOST")
	rabbitPort := os.Getenv("RABBIT_PORT")
	rabbitUrl := fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitUserName,rabbitPassword,rabbitHost,rabbitPort)
	return rabbitUrl
}

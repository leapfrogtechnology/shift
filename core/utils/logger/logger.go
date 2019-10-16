package logger

import (
	"log"

	"github.com/logrusorgru/aurora"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("Fatal\n%s: %s", aurora.Red(msg), aurora.Red(err))
	}
}

func LogError(err error, msg string) {
	log.Printf("%s: %s", aurora.Red(msg), aurora.Red(err))
}

func LogOutput(msg string) {
	log.Println(aurora.Cyan(msg))
}

func LogInfo(msg string) {
	log.Println(aurora.Blue(msg))
}

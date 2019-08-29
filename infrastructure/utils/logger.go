package utils

import "log"

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func LogError(err error, msg string) {
	log.Printf("%s: %s", msg, err)
}

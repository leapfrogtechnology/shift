package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func failOnError(err error, msg string) {
	if err != nil {
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

// GetSession returns Session for AWS.
func GetSession(profile string, region string) *session.Session {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           profile,
		Config: aws.Config{
			Region: aws.String(region),
		},
	})

	if err != nil {
		failOnError(err, "Could not connect to AWS.")
	}

	return sess
}

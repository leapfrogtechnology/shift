package main

import (
	"encoding/json"

	"github.com/leapfrogtechnology/shift/deployment/domain/project"
	"github.com/leapfrogtechnology/shift/deployment/internals/frontend"
	"github.com/leapfrogtechnology/shift/deployment/services/aws/s3"
	"github.com/leapfrogtechnology/shift/deployment/services/mq/deployment"
	"github.com/leapfrogtechnology/shift/deployment/services/mq/trigger"
	"github.com/leapfrogtechnology/shift/deployment/services/storage"
)

func deploy(msg []byte) {
	projectResponse := project.Response{}
	json.Unmarshal(msg, &projectResponse)

	buildData := frontend.BuildData{
		GitToken:     projectResponse.Deployment.GitToken,
		Platform:     projectResponse.Deployment.Platform,
		CloneURL:     projectResponse.Deployment.CloneURL,
		BuildCommand: projectResponse.Deployment.BuildCommand,
		DistFolder:   projectResponse.Deployment.DistFolder,
		AccessKey:    projectResponse.Deployment.AccessKey,
		SecretKey:    projectResponse.Deployment.SecretKey,
	}

	frontend.Build(buildData)

	s3.Deploy(s3.Data{
		AccessKey:  projectResponse.Deployment.AccessKey,
		SecretKey:  projectResponse.Deployment.SecretKey,
		Bucket:     projectResponse.Data.BucketName,
		URL:        projectResponse.Data.URL,
		DistFolder: projectResponse.Deployment.DistFolder,
	})

	storage.Save(projectResponse)
}

func triggerDeploy(msg []byte) {
	triggerRequest := project.TriggerRequest{}
	json.Unmarshal(msg, &triggerRequest)

}

func main() {
	deployment.Consume(deploy)
	trigger.Consume(triggerDeploy)
}

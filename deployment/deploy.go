package main

import (
	"encoding/json"
	"fmt"

	"github.com/leapfrogtechnology/shift/deployment/domain/project"
	"github.com/leapfrogtechnology/shift/deployment/internals/frontend"
	"github.com/leapfrogtechnology/shift/deployment/services/aws/s3"
	"github.com/leapfrogtechnology/shift/deployment/services/mq/deployment"
	"github.com/leapfrogtechnology/shift/deployment/services/mq/trigger"
	"github.com/leapfrogtechnology/shift/deployment/services/storage"
	"github.com/leapfrogtechnology/shift/deployment/utils/slack"
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
		Bucket:     projectResponse.Data.BucketName.Value,
		URL:        projectResponse.Data.FrontendWebURL.Value,
		DistFolder: projectResponse.Deployment.DistFolder,
	})

	storage.Save(projectResponse)

	slack.Notify(
		projectResponse.Deployment.SlackURL,
		fmt.Sprintf(
			"Successfull deploy of *%s* *%s* \n %s",
			projectResponse.ProjectName,
			projectResponse.Deployment.Name,
			projectResponse.Data.FrontendWebURL.Value),
		"#04EBB8")
}

func triggerDeploy(msg []byte) {
	triggerRequest := project.TriggerRequest{}
	json.Unmarshal(msg, &triggerRequest)

	jsonData := storage.Read()

	deploymentData := jsonData[triggerRequest.Project][triggerRequest.Deployment]

	fmt.Println(deploymentData)

	if _, ok := jsonData[triggerRequest.Project][triggerRequest.Deployment]; ok {
		slack.Notify(
			deploymentData.Deployment.SlackURL,
			fmt.Sprintf(
				"*There is a new deployment in progress.* \n Project: `%s` \n Deployment: `%s` \n Started by: `%s`",
				deploymentData.ProjectName,
				deploymentData.Deployment.Name,
				triggerRequest.User),
			"#1CA7FB")

		deploymentDataJSON, _ := json.Marshal(deploymentData)

		deployment.Publish(deploymentDataJSON)
	} else {
		fmt.Println("Deployment " + triggerRequest.Deployment + " for Project " + triggerRequest.Project + " not found")
	}
}

func main() {
	go trigger.Consume(triggerDeploy)
	deployment.Consume(deploy)
}

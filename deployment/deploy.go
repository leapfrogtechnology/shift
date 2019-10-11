package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/leapfrogtechnology/shift/deployment/domain/project"
	"github.com/leapfrogtechnology/shift/deployment/internals/backend"
	"github.com/leapfrogtechnology/shift/deployment/internals/frontend"
	"github.com/leapfrogtechnology/shift/deployment/services/aws/s3"
	"github.com/leapfrogtechnology/shift/deployment/services/mq/deployment"
	"github.com/leapfrogtechnology/shift/deployment/services/mq/trigger"
	"github.com/leapfrogtechnology/shift/deployment/services/storage"
	"github.com/leapfrogtechnology/shift/deployment/utils/slack"
)

func deployFrontend(projectResponse project.Response) {
	buildData := frontend.BuildData{
		GitToken:     projectResponse.Deployment.GitToken,
		Platform:     projectResponse.Deployment.Platform,
		CloneURL:     projectResponse.Deployment.CloneURL,
		BuildCommand: projectResponse.Deployment.BuildCommand,
		DistFolder:   projectResponse.Deployment.DistFolder,
		AccessKey:    projectResponse.Deployment.AccessKey,
		SecretKey:    projectResponse.Deployment.SecretKey,
	}
	error := frontend.Build(buildData)

	if error != nil {
		slack.Notify(
			projectResponse.Deployment.SlackURL,
			fmt.Sprintf(
				"Error: Deployment of *%s* *%s* failed. \n %s",
				projectResponse.ProjectName,
				projectResponse.Deployment.Name,
				"```"+error.Error()+"```"),
			"#FF6871")

		return
	}

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

func deploy(msg []byte) {
	projectResponse := project.Response{}
	json.Unmarshal(msg, &projectResponse)
	if strings.EqualFold(projectResponse.Deployment.Type, "frontend") {
		deployFrontend(projectResponse)
	} else if strings.EqualFold(projectResponse.Deployment.Type, "backend") {
		out, err := backend.Deploy(msg)
		if err != nil {
			slack.Notify(
				projectResponse.Deployment.SlackURL,
				fmt.Sprintf(
					"Error: Deployment of *%s* *%s* failed. \n %s",
					projectResponse.ProjectName,
					projectResponse.Deployment.Name,
					"```"+err.Error()+"```"),
				"#FF6871")

			return
		}
		storage.Save(projectResponse)
		slack.Notify(
			projectResponse.Deployment.SlackURL,
			fmt.Sprintf(
				"Successfull deploy of *%s* *%s* \n %s",
				projectResponse.ProjectName,
				projectResponse.Deployment.Name,
				out),
			"#04EBB8")
	} else {
		slack.Notify(
			projectResponse.Deployment.SlackURL,
			fmt.Sprintf(
				"Error: Deployment of *%s* *%s* failed. \n %s",
				projectResponse.ProjectName,
				projectResponse.Deployment.Name,
				"Unknown Deployment Type"),
			"#FF6871")
	}
}

func triggerDeploy(msg []byte) {
	triggerRequest := project.TriggerRequest{}
	json.Unmarshal(msg, &triggerRequest)

	jsonData := storage.Read()

	deploymentData := jsonData[triggerRequest.Project][triggerRequest.Deployment]

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

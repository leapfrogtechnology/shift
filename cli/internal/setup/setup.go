package setup

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/AlecAivazis/survey/v2"
)

type projectDetails struct {
	ProjectName string
}

type deploymentDetails struct {
	DeploymentName string
	CloudProvider  string
	DeploymentType string
}

type frontendBuildInformation struct {
	BuildCommand string
	DistFolder   string
}

type backendBuildInformation struct {
	Port            string
	HealthCheckPath string
	DockerfilePath  string
}

type deployment struct {
	Name            string `json:"name"`
	Platform        string `json:"platform"`
	Type            string `json:"type"`
	BuildCommand    string `json:"buildCommand"`
	DistFolder      string `json:"distFolder"`
	Port            string `json:"port"`
	HealthCheckPath string `json:"healthCheckPath"`
	SlackURL        string `json:"slackURL"`
	DockerFilePath  string `json:"dockerFilePath"`
}

// Project defines the overall structure for a project deployment.
type Project struct {
	ProjectName string     `json:"projectName"`
	Deployment  deployment `json:"deployment"`
}

func askProjectDetails() *projectDetails {
	questions := []*survey.Question{
		{
			Name: "projectName",
			Prompt: &survey.Input{
				Message: "Project Name",
			},
		},
	}

	answers := &projectDetails{}
	err := survey.Ask(questions, answers)

	if err != nil {
		fmt.Println(err)
	}

	return answers
}

func askDeploymentDetails() *deploymentDetails {
	questions := []*survey.Question{
		{
			Name: "deploymentName",
			Prompt: &survey.Input{
				Message: "Deployment Name:",
			},
		},
		{
			Name: "cloudProvider",
			Prompt: &survey.Select{
				Message: "Choose Cloud Provider:",
				Options: []string{"AWS", "Azure", "GCP"},
			},
		},
		{
			Name: "deploymentType",
			Prompt: &survey.Select{
				Message: "Choose Deployment Type:",
				Options: []string{"Frontend", "Backend"},
			},
		},
	}

	answers := &deploymentDetails{}
	err := survey.Ask(questions, answers)

	if err != nil {
		fmt.Println(err)
	}

	return answers
}

func askFrontendBuildInformation() *frontendBuildInformation {
	questions := []*survey.Question{
		{
			Name: "buildCommand",
			Prompt: &survey.Input{
				Message: "Build Command: ",
			},
		},
		{
			Name: "distFolder",
			Prompt: &survey.Input{
				Message: "Distribution Folder: ",
			},
		},
	}

	answers := &frontendBuildInformation{}
	err := survey.Ask(questions, answers)

	if err != nil {
		fmt.Println(err)
	}

	return answers
}

func askBackendBuildInformation() *backendBuildInformation {
	questions := []*survey.Question{
		{
			Name: "port",
			Prompt: &survey.Input{
				Message: "Application Port: ",
			},
		},
		{
			Name: "healthCheckPath",
			Prompt: &survey.Input{
				Message: "Healthcheck Path (eg: '/api'): ",
			},
		},
		{
			Name: "dockerfilePath",
			Prompt: &survey.Input{
				Message: "Dockerfile Path (eg: './'): ",
			},
		},
	}

	answers := &backendBuildInformation{}

	err := survey.Ask(questions, answers)

	if err != nil {
		fmt.Println(err)
	}

	return answers
}

func askSlackEndpoint() string {
	slackEndpoint := ""
	prompt := &survey.Input{
		Message: "Slack webhook URL: ",
	}
	survey.AskOne(prompt, &slackEndpoint)

	return slackEndpoint
}

// Run initializes setup for shift projects.
func Run() {
	projectDetails := askProjectDetails()
	deploymentDetails := askDeploymentDetails()

	frontendBuildInformation := &frontendBuildInformation{}

	backendBuildInformation := &backendBuildInformation{}

	if deploymentDetails.DeploymentType == "Frontend" {
		frontendBuildInformation = askFrontendBuildInformation()
	} else if deploymentDetails.DeploymentType == "Backend" {
		backendBuildInformation = askBackendBuildInformation()
	}

	slackEndpoint := askSlackEndpoint()

	projectRequest := Project{
		ProjectName: projectDetails.ProjectName,
		Deployment: deployment{
			Name:            deploymentDetails.DeploymentName,
			Platform:        deploymentDetails.CloudProvider,
			Type:            deploymentDetails.DeploymentType,
			BuildCommand:    frontendBuildInformation.BuildCommand,
			DistFolder:      frontendBuildInformation.DistFolder,
			Port:            backendBuildInformation.Port,
			HealthCheckPath: backendBuildInformation.HealthCheckPath,
			SlackURL:        slackEndpoint,
			DockerFilePath:  backendBuildInformation.DockerfilePath,
		},
	}

	// projectRequestJSON, _ := json.Marshal(projectRequest)

	// fmt.Println(string(projectRequestJSON))

	jsonData, _ := json.MarshalIndent(projectRequest, "", " ")

	currentDir, _ := os.Getwd()
	fileName := currentDir + "/shift.json"

	_ = ioutil.WriteFile(fileName, jsonData, 0644)

	// 2. Run infrastructre code here and save to JSON again

	// 3. Deploy to infrastructure code here.
}

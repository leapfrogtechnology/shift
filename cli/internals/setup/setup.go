package setup

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/leapfrogtechnology/shift/cli/internals/config"
	"github.com/leapfrogtechnology/shift/core/structs"
)

type projectDetails struct {
	ProjectName string
}

type deploymentDetails struct {
	CloudProvider  string
	Profile        string
	Region         string
	DeploymentType string
	Environment    string
}

type frontendBuildInformation struct {
	DistFolder string
}

type backendBuildInformation struct {
	Port            string
	HealthCheckPath string
	DockerfilePath  string
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
			Name: "cloudProvider",
			Prompt: &survey.Select{
				Message: "Choose Cloud Provider:",
				Options: []string{"AWS", "Azure", "GCP"},
			},
		},
		{
			Name: "Profile",
			Prompt: &survey.Select{
				Message: "Chose Aws Profile:",
				Options: config.GetProfiles(),
			},
		},
		{
			Name: "Region",
			Prompt: &survey.Select{
				Message: "Region:",
				Options: config.GetRegions(),
			},
		},
		{
			Name: "deploymentType",
			Prompt: &survey.Select{
				Message: "Choose Deployment Type: ",
				Options: []string{"Frontend", "Backend"},
			},
		},
		{
			Name: "environment",
			Prompt: &survey.Input{
				Message: "Environment name: ",
			},
		},
	}

	answers := &deploymentDetails{}
	err := survey.Ask(questions, answers)
	answers.Region = config.GetRegionCode(answers.Region)

	if err != nil {
		fmt.Println(err)
	}

	return answers
}

func askFrontendBuildInformation() *frontendBuildInformation {
	questions := []*survey.Question{
		{
			Name: "distFolder",
			Prompt: &survey.Input{
				Message: "Build Directory: ",
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

	projectRequest := structs.Project{
		Name:            projectDetails.ProjectName,
		Platform:        deploymentDetails.CloudProvider,
		Profile:         deploymentDetails.Profile,
		Region:          deploymentDetails.Region,
		Type:            deploymentDetails.DeploymentType,
		DistDir:         frontendBuildInformation.DistFolder,
		Port:            backendBuildInformation.Port,
		HealthCheckPath: backendBuildInformation.HealthCheckPath,
		SlackURL:        slackEndpoint,
		DockerFilePath:  backendBuildInformation.DockerfilePath,
		Env: map[string]structs.Env{
			deploymentDetails.Environment: structs.Env{},
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

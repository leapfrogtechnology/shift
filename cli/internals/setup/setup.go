package setup

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/leapfrogtechnology/shift/cli/internals/deploy"
	"github.com/leapfrogtechnology/shift/core/services/platforms/aws"
	"github.com/leapfrogtechnology/shift/core/services/storage"
	"github.com/leapfrogtechnology/shift/core/structs"
	"github.com/leapfrogtechnology/shift/infrastructure/internals/initialize"
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
				Options: aws.GetProfiles(),
			},
		},
		{
			Name: "Region",
			Prompt: &survey.Select{
				Message: "Region:",
				Options: aws.GetRegions(),
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
	answers.Region = aws.GetRegionCode(answers.Region)

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

	// 1. Save project details to shift.json.
	storage.Save(projectRequest)

	// 2. Run infrastructre code and save to JSON again with updated information.
	infraInfo := initialize.Run(projectRequest, deploymentDetails.Environment)

	// 3. Save Infrastructure details to shift.json.
	projectRequest.Env[deploymentDetails.Environment] = structs.Env{
		Bucket: infraInfo.Bucket,
	}

	storage.Save(projectRequest)

	// 3. Deploy to created infrastructure.
	deploy.Run(deploymentDetails.Environment)
}

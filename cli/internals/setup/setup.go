package setup

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/leapfrogtechnology/shift/cli/services/mq"
	"github.com/leapfrogtechnology/shift/cli/utils/github"
	"github.com/leapfrogtechnology/shift/cli/utils/spinner"
)

type projectDetails struct {
	ProjectName string
}

type deploymentDetails struct {
	DeploymentName string
	CloudProvider  string
	AccessKey      string
	SecretKey      string
	DeploymentType string
	GitProvider    string
}

type deployment struct {
	Name         string `json:"name"`
	Platform     string `json:"platform"`
	AccessKey    string `json:"accessKey"`
	SecretKey    string `json:"secretKey"`
	Type         string `json:"type"`
	GitProvider  string `json:"gitProvider"`
	GitToken     string `json:"gitToken"`
	CloneURL     string `json:"cloneURL"`
	BuildCommand string `json:"buildCommand"`
	DistFolder   string `json:"distFolder"`
}

type projectRequest struct {
	ProjectName string     `json:"projectName"`
	Deployment  deployment `json:"deployment"`
}

type buildInformation struct {
	BuildCommand string
	DistFolder   string
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
			Name: "accessKey",
			Prompt: &survey.Input{
				Message: "Access Key:",
			},
		},
		{
			Name: "secretKey",
			Prompt: &survey.Input{
				Message: "Secret Key:",
			},
		},
		{
			Name: "deploymentType",
			Prompt: &survey.Select{
				Message: "Choose Deployment Type:",
				Options: []string{"Frontend", "Backend"},
			},
		},
		{
			Name: "gitProvider",
			Prompt: &survey.Select{
				Message: "Choose Git Provider",
				Options: []string{"Github", "Gitlab", "BitBucket"},
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

func askGitCredentials(gitProvider string) *github.GitCredentials {
	fmt.Println("\n Connect " + gitProvider)
	questions := []*survey.Question{
		{
			Name: "username",
			Prompt: &survey.Input{
				Message: "Username",
			},
		},
		{
			Name: "password",
			Prompt: &survey.Password{
				Message: "Password",
			},
		},
	}

	answers := &github.GitCredentials{}
	err := survey.Ask(questions, answers)

	if err != nil {
		fmt.Println(err)
	}

	return answers
}

func chooseOrganization(personalToken string) string {
	spinner.Start("Fetching your organizations...")
	user, _ := github.FetchUser(personalToken)
	organizations, _ := github.FetchOrganizations(personalToken)
	spinner.Stop()

	questions := &survey.Select{
		Message: "Choose user/organization:",
		Options: append(organizations, user+" (User)"),
	}

	org := ""
	err := survey.AskOne(questions, &org)

	if err != nil {
		fmt.Println(err)
	}

	return org
}

func chooseRepo(personalToken string, organization string) (string, string) {
	repos := []string{}
	repoURL := map[string]string{}

	spinner.Start("Fetching your repositories...")
	if strings.Contains(organization, "(User)") {
		repos, repoURL, _ = github.FetchUserRepos(personalToken)
	} else {
		repos, repoURL, _ = github.FetchOrgRepos(personalToken, organization)
	}
	spinner.Stop()

	questions := &survey.Select{
		Message: "Choose Repository:",
		Options: repos,
	}

	org := ""
	err := survey.AskOne(questions, &org)

	if err != nil {
		fmt.Println(err)
	}

	return org, repoURL[org]
}

func askBuildInformation() *buildInformation {
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

	answers := &buildInformation{}
	err := survey.Ask(questions, answers)

	if err != nil {
		fmt.Println(err)
	}

	return answers
}

// Run initializes setup for shift projects.
func Run() {
	projectDetails := askProjectDetails()
	deploymentDetails := askDeploymentDetails()
	gitCredentials := askGitCredentials(deploymentDetails.GitProvider)

	spinner.Start("Connecting to Github...")
	personalToken, _ := github.CreatePersonalToken(gitCredentials)
	spinner.Stop()

	organization := chooseOrganization(personalToken)
	_, repoURL := chooseRepo(personalToken, organization)

	buildInformation := askBuildInformation()

	projectRequest := projectRequest{
		ProjectName: projectDetails.ProjectName,
		Deployment: deployment{
			Name:         deploymentDetails.DeploymentName,
			Platform:     deploymentDetails.CloudProvider,
			AccessKey:    deploymentDetails.AccessKey,
			SecretKey:    deploymentDetails.SecretKey,
			Type:         deploymentDetails.DeploymentType,
			GitProvider:  deploymentDetails.GitProvider,
			GitToken:     personalToken,
			CloneURL:     repoURL,
			BuildCommand: buildInformation.BuildCommand,
			DistFolder:   buildInformation.DistFolder,
		},
	}

	projectRequestJSON, _ := json.Marshal(projectRequest)

	fmt.Println(string(projectRequestJSON))

	mq.Publish(projectRequestJSON)
}
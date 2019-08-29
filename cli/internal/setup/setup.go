package setup

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
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
	spinner.Start()
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

func chooseRepo(personalToken string, organization string) string {
	repos := []string{}

	spinner.Start()
	if strings.Contains(organization, "(User)") {
		repos, _ = github.FetchUserRepos(personalToken)
	} else {
		repos, _ = github.FetchOrgRepos(personalToken, organization)
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

	return org
}

// Run initializes setup for shift projects.
func Run() {
	projectDetails := askProjectDetails()
	deploymentDetails := askDeploymentDetails()
	gitCredentials := askGitCredentials(deploymentDetails.GitProvider)

	spinner.Start()
	personalToken, _ := github.CreatePersonalToken(gitCredentials)
	spinner.Stop()

	organization := chooseOrganization(personalToken)
	repo := chooseRepo(personalToken, organization)

	fmt.Print("ProjectDetails: ")
	fmt.Println(projectDetails)

	fmt.Print("DeploymentDetails: ")
	fmt.Println(deploymentDetails)

	fmt.Print("Access Token: ")
	fmt.Println(personalToken)

	fmt.Print("Organization: ")
	fmt.Println(organization)

	fmt.Print("Repository: ")
	fmt.Println(repo)
}

package setup

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

type projectDetails struct {
	ProjectName string
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

func Run() {
	projectDetails := askProjectDetails()

	fmt.Println(projectDetails)

}

package template

import (
	"io/ioutil"
	"os"

	"github.com/flosch/pongo2"
	"github.com/leapfrogtechnology/shift/core/structs"
	"github.com/leapfrogtechnology/shift/core/utils/logger"
)

type infrastructure struct {
	Client      structs.Project
	Token       string
	Environment string
}

// GenerateFrontendTemplate generates template for the frontend infrastructure.
func GenerateFrontendTemplate(project structs.Project, terraformPath string, environment string) {
	token := os.Getenv("TERRAFORM_TOKEN")
	infrastructure := infrastructure{project, token, environment}

	tpl, err := pongo2.FromString(FrontendTemplate)

	logger.FailOnError(err, "Failed to parse string")

	out, err := tpl.Execute(pongo2.Context{"info": infrastructure})

	logger.FailOnError(err, "Failed to parse string")

	terraformFileName := terraformPath + "/infrastructure.tf"
	err = os.MkdirAll(terraformPath, 0700)

	logger.FailOnError(err, "Failed to create terraform directory")

	err = ioutil.WriteFile(terraformFileName, []byte(out), 0600)

	logger.FailOnError(err, "Failed to create terraform directory")
}

// GenerateBackendTemplate generates template for the backend infrastructure.
func GenerateBackendTemplate(project structs.Project, terraformPath string, environment string) {
	token := os.Getenv("TERRAFORM_TOKEN")
	infrastructure := infrastructure{project, token, environment}

	tpl, err := pongo2.FromString(BackendTemplate)

	logger.FailOnError(err, "Failed to parse string")

	out, err := tpl.Execute(pongo2.Context{"info": infrastructure})

	logger.FailOnError(err, "Failed to parse string")

	terraformFileName := terraformPath + "/infrastructure.tf"
	err = os.MkdirAll(terraformPath, 0700)

	logger.FailOnError(err, "Failed to create terraform directory")

	err = ioutil.WriteFile(terraformFileName, []byte(out), 0600)

	logger.FailOnError(err, "Failed to create terraform template")
}

// GenerateContainerTemplate generates template to run initial container.
func GenerateContainerTemplate(containerTemplateFile string) {
	err := ioutil.WriteFile(containerTemplateFile, []byte(ContainerTemplate), 0600)

	logger.FailOnError(err, "Failed to create container template")
}

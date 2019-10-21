package template

import (
	"fmt"
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

	fmt.Println(out)
}

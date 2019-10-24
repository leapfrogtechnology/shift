package destroy

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/leapfrogtechnology/shift/core/services/storage"
	"github.com/leapfrogtechnology/shift/core/utils/file"
	"github.com/leapfrogtechnology/shift/core/utils/system/exit"
	"github.com/leapfrogtechnology/shift/infrastructure/internals/terraform"
)

func askConformation(environment, projectName string) string {
	conformation := ""
	prompt := &survey.Input{
		Message: "Are you sure you want to destroy " + environment + " environment from " + projectName + " ?(Y/N): ",
	}
	survey.AskOne(prompt, &conformation)

	return conformation
}

// Run initializes destruction of infrastructure
func Run(environment string) {
	project := storage.Read()
	_, env := project.Env[environment]

	if !env {
		const message = "Unknown Environment type "
		exit.Error(errors.New(message+"'"+environment+"'"), "Error")
	}

	conformation := askConformation(environment, project.Name)

	if strings.EqualFold(conformation, "Y") || strings.EqualFold(conformation, "yes") {

		workspaceRoot := "/tmp"
		workspaceDir := filepath.Join(workspaceRoot, project.Name, project.Type, environment)
		terraformFile := workspaceDir + "/infrastructure.tf"

		exists := file.Exists(terraformFile)

		if exists {
			terraform.DestroyInfrastructure(workspaceDir)
		} else {
			terraform.MakeTempAndDestroy(project, environment, workspaceDir)
		}

	} else {
		const message = "Denied by user"
		exit.Error(errors.New(message), "Cancelled")
	}
}

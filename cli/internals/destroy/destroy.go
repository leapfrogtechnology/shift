package destroy

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/leapfrogtechnology/shift/core/services/slack"
	"github.com/leapfrogtechnology/shift/core/services/storage"
	"github.com/leapfrogtechnology/shift/core/utils/file"
	"github.com/leapfrogtechnology/shift/core/utils/system"
	"github.com/leapfrogtechnology/shift/core/utils/system/exit"
	"github.com/leapfrogtechnology/shift/infrastructure/internals/terraform"
)

// Run initializes new environment.
func Run(environment string) {
	project := storage.Read()
	_, ok := project.Env[environment]

	if !ok {
		exit.Error(errors.New("Unknown Environment type "+"'"+environment+"'"), "Error")
	}

	workspaceRoot := "/tmp"

	workspaceDir := filepath.Join(workspaceRoot, project.Name, project.Type, environment)

	terraformFile := workspaceDir + "/infrastructure.tf"

	slack.Notify(project.SlackURL,
		fmt.Sprintf("*There is Infrastructure Destroy in progress* \n Project: `%s` \n Environment: `%s` \n Started by: `%s`",
			project.Name, environment, system.CurrentUser()),
		"#1CA7FB")

	isExist := file.IsExist(terraformFile)

	if isExist {
		err := terraform.DestroyInfrastructure(workspaceDir)

		slack.Notify(project.SlackURL, fmt.Sprintf("Error :\n %s", err.Error()), "#04EBB8")
		os.Exit(1)

	} else {
		err := terraform.MakeTempAndDestroy(project, environment, workspaceDir)
		slack.Notify(project.SlackURL, fmt.Sprintf("Error :\n %s", err.Error()), "#04EBB8")
		os.Exit(1)
	}

	slack.Notify(project.SlackURL, fmt.Sprintf(" 🎉 🎉 🎉Successfully destroyed *%s* env of *%s* from *%s*. 🎉 🎉 🎉", environment, project.Name, project.Platform), "#04EBB8")
}

package destroy

import (
	"errors"
	"path/filepath"

	"github.com/leapfrogtechnology/shift/core/services/storage"
	"github.com/leapfrogtechnology/shift/core/utils/file"
	"github.com/leapfrogtechnology/shift/core/utils/system/exit"
	"github.com/leapfrogtechnology/shift/infrastructure/internals/terraform"
)

// Run initializes destruction of infrastructure
func Run(environment string) {

	project := storage.Read()
	_, env := project.Env[environment]

	if !env {
		exit.Error(errors.New("Unknown Environment type "+"'"+environment+"'"), "Error")
	}

	workspaceRoot := "/tmp"

	workspaceDir := filepath.Join(workspaceRoot, project.Name, project.Type, environment)

	terraformFile := workspaceDir + "/infrastructure.tf"

	isExist := file.Exists(terraformFile)

	if isExist {
		terraform.DestroyInfrastructure(workspaceDir)

	} else {
		terraform.MakeTempAndDestroy(project, environment, workspaceDir)
	}
}

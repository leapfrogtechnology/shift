package infrastrucuture

import (
	"encoding/json"
	"github.com/leapfrogtechnology/shift/infrastructure/templates/providers/aws/frontend-architecture"
	"github.com/leapfrogtechnology/shift/infrastructure/utils"
	"path/filepath"

	"github.com/leapfrogtechnology/shift/infrastructure/utils"
)

func InitializeFrontend(infrastructureArgs string) error{

	var frontendArgs utils.FrontendInfrastructureVariables
	err := json.Unmarshal([]byte(infrastructureArgs), &frontendArgs)
	if err != nil {
		utils.LogError(err, "Error Parsing Body")
		return err
	}
	workspaceDir := filepath.Join("/tmp", frontendArgs.CLIENT_NAME)

	err = utils.GenerateFrontendTemplateFile(frontend_architecture.InfrastructureTemplate, frontendArgs, workspaceDir)
	if err != nil {
		utils.LogError(err, "Cannot Generate Template")
		return err
	}
	utils.RunInfrastructureChanges(workspaceDir)
	return nil
}

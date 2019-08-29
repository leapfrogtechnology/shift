package infrastrucuture

import (
	"encoding/json"
	"github.com/leapfrogtechnology/shift/infrastructure/templates/providers/aws/frontend-architecture"
	"github.com/leapfrogtechnology/shift/infrastructure/utils"
	"path/filepath"

)

func InitializeFrontend(infrastructureArgs []byte) (string, error){

	var frontendArgs utils.FrontendInfrastructureVariables
	utils.LogInfo("Gathering Info")
	err := json.Unmarshal(infrastructureArgs, &frontendArgs)
	if err != nil {
		utils.LogError(err, "Error Parsing Body")
		return "", err
	}
	workspaceDir := filepath.Join("/tmp", frontendArgs.CLIENT_NAME)
	utils.LogInfo("Generating Template")
	err = utils.GenerateFrontendTemplateFile(frontend_architecture.InfrastructureTemplate, frontendArgs, workspaceDir)
	if err != nil {
		utils.LogError(err, "Cannot Generate Template")
		return "", err
	}
	utils.LogInfo("Running Infrastructure Changes")
	infrastructureInfo, err := utils.RunInfrastructureChanges(workspaceDir)
	if err != nil {
		utils.LogError(err, "Cannot Run Changes")
		return "", err
	}
	return infrastructureInfo, err
}

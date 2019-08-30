package infrastrucuture

import (
	"encoding/json"
	"github.com/leapfrogtechnology/shift/infrastructure/templates/providers/aws/frontend-architecture"
	"github.com/leapfrogtechnology/shift/infrastructure/utils"
	"path/filepath"

)

func InitializeFrontend(ClientArgs []byte) (string, error){

	var clientArgs utils.Client
	utils.LogInfo("Gathering Info")
	err := json.Unmarshal(ClientArgs, &clientArgs)
	if err != nil {
		utils.LogError(err, "Error Parsing Body")
		return "", err
	}
	workspaceDir := filepath.Join("/tmp", clientArgs.Deployment.Name)
	utils.LogInfo("Generating Template")
	err = utils.GenerateFrontendTemplateFile(frontend_architecture.InfrastructureTemplate, clientArgs, workspaceDir)
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

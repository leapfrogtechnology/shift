package infrastrucuture

import (
	"encoding/json"
	"github.com/leapfrogtechnology/shift/core/utils/logger"
	"github.com/leapfrogtechnology/shift/infrastructure/internal/terraform"
	"github.com/leapfrogtechnology/shift/infrastructure/internal/terraform/templates/providers/aws/frontend-architecture"
	"github.com/leapfrogtechnology/shift/infrastructure/utils"
	"path/filepath"
)

func InitializeFrontend(ClientArgs []byte) (string, error) {

	var clientArgs utils.Client
	logger.LogInfo("Gathering Info")
	err := json.Unmarshal(ClientArgs, &clientArgs)
	if err != nil {
		logger.LogError(err, "Error Parsing Body")
		return "", err
	}
	workspaceDir := filepath.Join("/tmp", clientArgs.Deployment.Name)
	logger.LogInfo("Generating Template")
	err = utils.GenerateFrontendTemplateFile(frontend_architecture.InfrastructureTemplate, clientArgs, workspaceDir)
	if err != nil {
		logger.LogError(err, "Cannot Generate Template")
		return "", err
	}
	logger.LogInfo("Running Infrastructure Changes")
	bucketName, url, err := terraform.RunFrontendInfrastructureChanges(workspaceDir)
	if err != nil {
		logger.LogError(err, "Cannot Run Changes")
		return "", err
	}
	result := utils.FrontendResult{
		Project:    clientArgs.Project,
		Deployment: clientArgs.Deployment,
		Data: utils.Frontend{
			BucketName: bucketName,
			Url:        url,
		},
	}
	out, err := json.Marshal(result)
	if err != nil {
		logger.LogError(err, "Error Marshalling output")
		return "", err
	}
	return string(out), err
}

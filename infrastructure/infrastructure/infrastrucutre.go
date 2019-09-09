package infrastrucuture

import (
	"encoding/json"
	"github.com/leapfrogtechnology/shift/core/utils/logger"
	"github.com/leapfrogtechnology/shift/infrastructure/internal/terraform"
	"github.com/leapfrogtechnology/shift/infrastructure/internal/terraform/templates/providers/aws/frontend-architecture"
	"github.com/leapfrogtechnology/shift/infrastructure/utils"
	"path/filepath"
)

type terraformOutput struct {
	Sensitive bool   `json:"sensitive"`
	Type      string `json:"type"`
	Value     string `json:"value"`
}
type FrontendTerraformOutput struct {
	BucketName     terraformOutput `json:"bucket_name"`
	FrontendWebUrl terraformOutput `json:"frontend_web_url"`
}

func InitializeFrontend(ClientArgs []byte) (string, error) {
	workspaceRoot := "/tmp"

	var clientArgs utils.Client
	var frontendTerraformOutput FrontendTerraformOutput

	logger.LogInfo("Gathering Info")

	err := json.Unmarshal(ClientArgs, &clientArgs)
	if err != nil {
		logger.LogError(err, "Error Parsing Body")
		return "", err
	}
	workspaceDir := filepath.Join(workspaceRoot, clientArgs.Deployment.Name)
	logger.LogInfo("Generating Template")
	err = utils.GenerateFrontendTemplateFile(frontend_architecture.InfrastructureTemplate, clientArgs, workspaceDir)
	if err != nil {
		logger.LogError(err, "Cannot Generate Template")
		return "", err
	}
	logger.LogInfo("Running Infrastructure Changes")
	terraformOutput, err := terraform.RunInfrastructureChanges(workspaceDir)
	if err != nil {
		logger.LogError(err, "Cannot Run Changes")
		return "", err
	}
	_ = json.Unmarshal([]byte(terraformOutput), &frontendTerraformOutput)
	result := utils.FrontendResult{
		Project:    clientArgs.Project,
		Deployment: clientArgs.Deployment,
		Data: utils.Frontend{
			BucketName: frontendTerraformOutput.BucketName.Value,
			Url:        frontendTerraformOutput.FrontendWebUrl.Value,
		},
	}
	out, err := json.Marshal(result)
	if err != nil {
		logger.LogError(err, "Error Marshalling output")
		return "", err
	}
	return string(out), err
}

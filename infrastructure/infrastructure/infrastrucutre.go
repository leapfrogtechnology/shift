package infrastrucuture

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/leapfrogtechnology/shift/core/utils/logger"
	"github.com/leapfrogtechnology/shift/infrastructure/internals/terraform"
	backendHaArchitecture "github.com/leapfrogtechnology/shift/infrastructure/internals/terraform/templates/providers/aws/backend-ha-architecture"
	frontendArchitecture "github.com/leapfrogtechnology/shift/infrastructure/internals/terraform/templates/providers/aws/frontend-architecture"
	"github.com/leapfrogtechnology/shift/infrastructure/utils"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func Initialize(ClientArgs []byte) (string, error) {
	var clientArgs utils.Client
	err := json.Unmarshal(ClientArgs, &clientArgs)
	if err != nil {
		logger.LogError(err, "Error Parsing Body")
		return "", err
	}
	if strings.EqualFold(clientArgs.Deployment.Type, "frontend") {
		out, err := InitializeFrontend(ClientArgs)
		return out, err
	} else if strings.EqualFold(clientArgs.Deployment.Type, "backend") {
		out, err := InitializeBackend(ClientArgs)
		return out, err
	} else {
		return "", errors.New(fmt.Sprintf("Unknown Deployment Type:\t%s", clientArgs.Deployment.Type))
	}

}

func InitializeFrontend(ClientArgs []byte) (string, error) {
	workspaceRoot := "/tmp"

	var clientArgs utils.Client
	var frontendTerraformOutput utils.FrontendTerraformOutput

	logger.LogInfo("Gathering Info")

	err := json.Unmarshal(ClientArgs, &clientArgs)
	if err != nil {
		logger.LogError(err, "Error Parsing Body")
		return "", err
	}
	workspaceDir := filepath.Join(workspaceRoot, clientArgs.Project, clientArgs.Deployment.Name, clientArgs.Deployment.Type)
	logger.LogInfo("Generating Template")
	err = utils.GenerateFrontendTemplateFile(frontendArchitecture.InfrastructureTemplate, clientArgs, workspaceDir)
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
		Data:       frontendTerraformOutput,
	}

	out, err := json.Marshal(result)
	if err != nil {
		logger.LogError(err, "Error Marshalling output")
		return "", err
	}
	logger.LogOutput(string(out))
	return string(out), err
}

func InitializeBackend(ClientArgs []byte) (string, error) {
	workspaceRoot := "/tmp"
	var clientArgs utils.Client
	var backendTerraformOutput utils.BackendTerraformOutput
	logger.LogInfo("Gathering Info")

	err := json.Unmarshal(ClientArgs, &clientArgs)
	if err != nil {
		logger.LogError(err, "Error Parsing Body")
		return "", err
	}
	workspaceDir := filepath.Join(workspaceRoot, clientArgs.Project, clientArgs.Deployment.Name, clientArgs.Deployment.Type)
	logger.LogInfo("Generating Template")
	err = utils.GenerateFrontendTemplateFile(backendHaArchitecture.InfrastructureTemplate, clientArgs, workspaceDir)
	containerTemplateFile := workspaceDir + "/sample.json.tpl"
	_ = ioutil.WriteFile(containerTemplateFile, []byte(backendHaArchitecture.ContainerTemplate), 0600)
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
	_ = json.Unmarshal([]byte(terraformOutput), &backendTerraformOutput)
	result := utils.BackendResult{
		Project:    clientArgs.Project,
		Deployment: clientArgs.Deployment,
		Data:       backendTerraformOutput,
	}
	out, err := json.Marshal(result)
	if err != nil {
		logger.LogError(err, "Error Marshalling output")
		return "", err
	}
	logger.LogOutput(string(out))

	return string(out), err
}

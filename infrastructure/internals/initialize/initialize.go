package initialize

import (
	"encoding/json"
	"errors"

	"strings"

	"path/filepath"

	"github.com/leapfrogtechnology/shift/core/utils/logger"
	"github.com/leapfrogtechnology/shift/infrastructure/utils"

	"github.com/leapfrogtechnology/shift/core/structs"
	"github.com/leapfrogtechnology/shift/infrastructure/internals/terraform"
	"github.com/leapfrogtechnology/shift/infrastructure/internals/terraform/templates/providers/aws/template"
	// backendHaArchitecture "github.com/leapfrogtechnology/shift/infrastructure/internals/terraform/templates/providers/aws/backend-ha-architecture"
)

type terraformOutput struct {
	Sensitive bool   `json:"sensitive"`
	Type      string `json:"type"`
	Value     string `json:"value"`
}

type frontendTerraformOutput struct {
	BucketName     terraformOutput `json:"bucketName"`
	FrontendWebURL terraformOutput `json:"appUrl"`
}

// // CreateBackend creates infrastructure for the backend.
// func CreateBackend(project structs.Project) (string, error) {
// 	workspaceRoot := "/tmp"
// 	var clientArgs utils.Client
// 	var backendTerraformOutput utils.BackendTerraformOutput
// 	logger.LogInfo("Gathering Info")

// 	err := json.Unmarshal(ClientArgs, &clientArgs)
// 	if err != nil {
// 		logger.LogError(err, "Error Parsing Body")
// 		return "", err
// 	}
// 	workspaceDir := filepath.Join(workspaceRoot, clientArgs.Project, clientArgs.Deployment.Name, clientArgs.Deployment.Type)
// 	logger.LogInfo("Generating Template")
// 	err = utils.GenerateFrontendTemplateFile(backendHaArchitecture.InfrastructureTemplate, clientArgs, workspaceDir)
// 	containerTemplateFile := workspaceDir + "/sample.json.tpl"
// 	_ = ioutil.WriteFile(containerTemplateFile, []byte(backendHaArchitecture.ContainerTemplate), 0600)
// 	if err != nil {
// 		logger.LogError(err, "Cannot Generate Template")
// 		return "", err
// 	}
// 	logger.LogInfo("Running Infrastructure Changes")
// 	terraformOutput, err := terraform.RunInfrastructureChanges(workspaceDir)
// 	if err != nil {
// 		logger.LogError(err, "Cannot Run Changes")
// 		return "", err
// 	}
// 	_ = json.Unmarshal([]byte(terraformOutput), &backendTerraformOutput)
// 	result := utils.BackendResult{
// 		Project:    clientArgs.Project,
// 		Deployment: clientArgs.Deployment,
// 		Data:       backendTerraformOutput,
// 	}
// 	out, err := json.Marshal(result)
// 	if err != nil {
// 		logger.LogError(err, "Error Marshalling output")
// 		return "", err
// 	}
// 	logger.LogOutput(string(out))

// 	return string(out), err
// }

// Run intitializes the infrastructure in the specified cloud provider.
func Run(project structs.Project, environment string) structs.Frontend {
	if !utils.CommandExists("terraform") {
		logger.FailOnError(errors.New("terraform does not exist"), "Please install terraform on your device")
	}

	infrastructureInfo := Initialize(project, environment)

	return infrastructureInfo
}

// Initialize creates infrastructure for frontend and backend according the given input.
func Initialize(project structs.Project, environment string) structs.Frontend {
	if strings.EqualFold(project.Type, "frontend") {
		out := CreateFrontend(project, environment)

		return out
	}
	// else if strings.EqualFold(project.Type, "backend") {
	// 	// out, err := CreateBackend(project)

	// 	return out
	// }

	logger.FailOnError(errors.New("Unknown Deployment Type"), project.Type)

	return structs.Frontend{}
}

// CreateFrontend creates infrastructure for frontend.
func CreateFrontend(project structs.Project, environment string) structs.Frontend {
	workspaceRoot := "/tmp"

	terraformOutput := frontendTerraformOutput{}

	workspaceDir := filepath.Join(workspaceRoot, project.Name, project.Type, environment)

	logger.LogInfo("Generating Template")
	template.GenerateFrontendTemplate(project, workspaceDir, environment)

	logger.LogInfo("Running Infrastructure Changes")
	workspaceName := project.Name + "-" + project.Type + "-" + environment
	output, err := terraform.RunInfrastructureChanges(workspaceDir, workspaceName)

	logger.FailOnError(err, "Cannot run changes")

	_ = json.Unmarshal([]byte(output), &terraformOutput)

	frontend := structs.Frontend{
		Bucket: terraformOutput.BucketName.Value,
		URL:    terraformOutput.FrontendWebURL.Value,
	}

	return frontend
}

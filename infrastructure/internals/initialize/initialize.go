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

type backendTerraformOutput struct {
	BackendClusterName         terraformOutput `json:"backendClusterName"`
	BackendContainerDefinition terraformOutput `json:"backendContainerDefinition"`
	BackendServiceID           terraformOutput `json:"backendServiceId"`
	BackendTaskDefinitionID    terraformOutput `json:"backendTaskDefinitionId"`
	BackendURL                 terraformOutput `json:"appUrl"`
	RepoURL                    terraformOutput `json:"repoUrl"`
}

// Run intitializes the infrastructure in the specified cloud provider.
func Run(project structs.Project, environment string) structs.Infrastructure {
	if !utils.CommandExists("terraform") {
		logger.FailOnError(errors.New("terraform does not exist"), "Please install terraform on your device")
	}

	infrastructureInfo := Initialize(project, environment)

	return infrastructureInfo
}

// Initialize creates infrastructure for frontend and backend according the given input.
func Initialize(project structs.Project, environment string) structs.Infrastructure {
	if strings.EqualFold(project.Type, "frontend") {
		out := CreateFrontend(project, environment)

		return out
	} else if strings.EqualFold(project.Type, "backend") {
		out := CreateBackend(project, environment)

		return out
	}

	logger.FailOnError(errors.New("Unknown Deployment Type"), project.Type)

	return structs.Infrastructure{}
}

// CreateFrontend creates infrastructure for frontend.
func CreateFrontend(project structs.Project, environment string) structs.Infrastructure {
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

	frontend := structs.Infrastructure{
		Bucket: terraformOutput.BucketName.Value,
		URL:    terraformOutput.FrontendWebURL.Value,
	}

	return frontend
}

// CreateBackend creates infrastructure for the backend.
func CreateBackend(project structs.Project, environment string) structs.Infrastructure {
	workspaceRoot := "/tmp"

	backendTerraformOutput := backendTerraformOutput{}

	logger.LogInfo("Gathering Info")

	workspaceDir := filepath.Join(workspaceRoot, project.Name, environment, project.Type)

	logger.LogInfo("Generating Template")
	template.GenerateBackendTemplate(project, workspaceDir, environment)

	containerTemplateFile := workspaceDir + "/sample.json.tpl"
	template.GenerateContainerTemplate(containerTemplateFile)

	logger.LogInfo("Running Infrastructure Changes")

	workspaceName := project.Name + "-" + project.Type + "-" + environment
	terraformOutput, err := terraform.RunInfrastructureChanges(workspaceDir, workspaceName)

	logger.FailOnError(err, "Cannot run changes")

	_ = json.Unmarshal([]byte(terraformOutput), &backendTerraformOutput)

	backend := structs.Infrastructure{
		Cluster:             backendTerraformOutput.BackendClusterName.Value,
		ContainerDefinition: backendTerraformOutput.BackendContainerDefinition.Value,
		ServiceID:           backendTerraformOutput.BackendServiceID.Value,
		TaskDefinitionID:    backendTerraformOutput.BackendTaskDefinitionID.Value,
		BackendURL:          backendTerraformOutput.BackendURL.Value,
		RepoURL:             backendTerraformOutput.RepoURL.Value,
	}

	return backend
}

package infrastrucuture

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/leapfrogtechnology/shift/core/utils/logger"
	"github.com/leapfrogtechnology/shift/infrastructure/internal/terraform"
	backendHaArchitecture "github.com/leapfrogtechnology/shift/infrastructure/internal/terraform/templates/providers/aws/backend-ha-architecture"
	frontendArchitecture "github.com/leapfrogtechnology/shift/infrastructure/internal/terraform/templates/providers/aws/frontend-architecture"
	"github.com/leapfrogtechnology/shift/infrastructure/utils"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
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

	// Deployment
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond) // Build our new spinner
	s.Prefix = "  "
	s.Suffix = "  Deploying"
	_ = s.Color("cyan", "bold")
	s.Start()
	region := "us-east-1"
	cloneUrl := "https://" + result.Deployment.GitToken + "@" + result.Deployment.CloneUrl[8:]
	command := fmt.Sprintf("rm -rf code && git clone %s code && $(AWS_ACCESS_KEY_ID=%s AWS_SECRET_ACCESS_KEY=%s aws ecr get-login --no-include-email --region %s) && docker build/%s code -t %s && docker push %s && terraform apply -var repo_url=%s -var port=%s --auto-approve", cloneUrl, result.Deployment.AccessKey, result.Deployment.SecretKey, region, result.Deployment.DockerFilePath,result.Data.RepoUrl.Value, result.Data.RepoUrl.Value, result.Data.RepoUrl.Value, result.Deployment.Port)
	logger.LogInfo(command)
	cmd := exec.Command("bash", "-c", command)
	cmd.Dir = workspaceDir
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	runError := cmd.Run()
	if runError != nil {
		logger.LogError(runError, stderr.String())
		s.Stop()
	} else {
		logger.LogOutput(stdout.String())
		s.Stop()
	}

	return string(out), err
}

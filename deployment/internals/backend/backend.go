package backend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/leapfrogtechnology/shift/core/utils/logger"
	backend_ha_architecture "github.com/leapfrogtechnology/shift/infrastructure/internals/terraform/templates/providers/aws/backend-ha-architecture"
	"github.com/leapfrogtechnology/shift/infrastructure/utils"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

func Deploy(resultJson []byte) (string ,error) {
	logger.LogInfo("Deploying")
	workspaceRoot := "/tmp"
	var result utils.BackendResult
	logger.LogInfo("Unmarshalling")
	_ = json.Unmarshal(resultJson, &result)
	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	result.Deployment.RepoName = result.Data.RepoUrl.Value + ":" + timestamp
	clientArgs := utils.Client{
		Project:    result.Project,
		Deployment: result.Deployment,
	}
	workspaceDir := filepath.Join(workspaceRoot, clientArgs.Project, clientArgs.Deployment.Name, clientArgs.Deployment.Type)
	logger.LogInfo("Generating Template")
	_ = utils.GenerateFrontendTemplateFile(backend_ha_architecture.InfrastructureTemplate, clientArgs, workspaceDir)
	logger.LogInfo("Generating Json")
	containerTemplateFile := workspaceDir + "/sample.json.tpl"
	_ = ioutil.WriteFile(containerTemplateFile, []byte(backend_ha_architecture.ContainerTemplate), 0600)

	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond) // Build our new spinner
	s.Prefix = "  "
	s.Suffix = "  Deploying"
	_ = s.Color("cyan", "bold")
	s.Start()
	region := "us-east-1"
	cloneUrl := "https://" + result.Deployment.GitToken + "@" + result.Deployment.CloneUrl[8:]
	logger.LogInfo(cloneUrl)
	command := fmt.Sprintf("rm -rf code && git clone %s code && $(AWS_ACCESS_KEY_ID=%s AWS_SECRET_ACCESS_KEY=%s aws ecr get-login --no-include-email --region %s) && docker build code/%s -t %s && docker push %s && terraform apply --auto-approve", cloneUrl, result.Deployment.AccessKey, result.Deployment.SecretKey, region, result.Deployment.DockerFilePath, result.Deployment.RepoName, result.Deployment.RepoName)
	logger.LogInfo(command)
	cmd := exec.Command("bash", "-c", command)
	cmd.Dir = workspaceDir
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		logger.LogError(err, stderr.String())
		s.Stop()
	} else {
		logger.LogOutput(stdout.String())
		s.Stop()
	}
	return result.Data.BackendUrl.Value , err
}

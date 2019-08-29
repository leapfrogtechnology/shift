package infrastructure

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/leapfrogtechnology/shift/infrastructure/utils"
)

func InitializeFrontend() {
	credentialsJsonFile, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	defer credentialsJsonFile.Close()
	byteValue, _ := ioutil.ReadAll(credentialsJsonFile)
	var frontendArgs utils.FrontendInfrastructureVariables
	err = json.Unmarshal(byteValue, &frontendArgs)
	if err != nil {
		panic(err)
	}
	workspaceDir := filepath.Join("/tmp", frontendArgs.CLIENT_NAME)
	utils.GenerateFrontendTemplateFile(frontend_architecture.InfrastructureTemplate, frontendArgs, workspaceDir)
	utils.RunInfrastructureChanges(workspaceDir)
}

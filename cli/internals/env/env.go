package env

import (
	"github.com/AlecAivazis/survey/v2"

	"github.com/leapfrogtechnology/shift/cli/internals/deploy"

	"github.com/leapfrogtechnology/shift/core/services/storage"
	"github.com/leapfrogtechnology/shift/core/structs"
	"github.com/leapfrogtechnology/shift/core/utils/logger"
	"github.com/leapfrogtechnology/shift/infrastructure/internals/initialize"
)

func askEnvironmentName() string {
	environment := ""
	prompt := &survey.Input{
		Message: "Environment Name (eg: dev): ",
	}

	survey.AskOne(prompt, &environment)

	return environment
}

// Run initializes new environment.
func Run() {
	project := storage.Read()
	environment := askEnvironmentName()

	infraInfo := initialize.Run(project, environment)

	project.Env[environment] = structs.Env{
		Bucket: infraInfo.Bucket,
	}

	storage.Save(project)

	deploy.Run(environment)

	logger.Info("Project Deployed At: " + infraInfo.URL)
	logger.Info("⏳ Stay put. It might take some time for the changes to be reflected. ⏳")
}

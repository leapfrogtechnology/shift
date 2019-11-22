package deploy

import (
	"errors"
	"fmt"
	"strings"

	"github.com/leapfrogtechnology/shift/core/services/slack"
	"github.com/leapfrogtechnology/shift/core/services/storage"
	"github.com/leapfrogtechnology/shift/core/utils/system"
	"github.com/leapfrogtechnology/shift/core/utils/system/exit"

	"github.com/leapfrogtechnology/shift/deployment/internals/backend"
	"github.com/leapfrogtechnology/shift/deployment/internals/frontend"
)

// Run starts deployment for the given environment
func Run(environment string) {
	project := storage.Read()

	_, ok := project.Env[environment]

	if !ok {
		exit.Error(errors.New("Unknown deployment type "+"'"+environment+"'"), "Error")
	}

	slack.Notify(project.SlackURL,
		fmt.Sprintf("*There is a new deployment in progress.* \n Project: `%s` \n Environment: `%s` \n Started by: `%s`",
			project.Name, environment, system.CurrentUser()),
		"#1CA7FB")

	if strings.EqualFold(project.Type, "frontend") {
		frontend.Deploy(project, environment)
	} else if strings.EqualFold(project.Type, "backend") {
		backend.Deploy(project, environment)
	}

	slack.Notify(project.SlackURL, fmt.Sprintf("*%s* succesfully deployed to *%s*. 🎉 🎉 🎉", project.Name, environment), "#04EBB8")
}

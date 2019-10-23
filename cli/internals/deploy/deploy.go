package deploy

import (
	"errors"
	"strings"

	"github.com/leapfrogtechnology/shift/core/services/storage"
	"github.com/leapfrogtechnology/shift/core/utils/system/exit"
	"github.com/leapfrogtechnology/shift/deployment/internals/frontend"
)

// Run starts deployment for the given environment
func Run(environment string) {
	project := storage.Read()

	_, ok := project.Env[environment]

	if !ok {
		exit.Error(errors.New("Unknown deployment type "+"'"+environment+"'"), "Error")
	}

	if strings.EqualFold(project.Type, "frontend") {
		frontend.Deploy(project, environment)
	}
}

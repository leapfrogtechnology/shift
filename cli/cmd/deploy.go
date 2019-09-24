package cmd

import (
	"github.com/leapfrogtechnology/shift/cli/internals/deploy"
)

// Deploy triggers deployment for provided project.
func Deploy(project string, deployment string) {
	deploy.Run(project, deployment)
}
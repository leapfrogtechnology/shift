package cmd

import (
	"github.com/leapfrogtechnology/shift/cli/internals/deploy"
)

// Deploy triggers deployment for the given environment.
func Deploy(environment string) {
	deploy.Run(environment)
}

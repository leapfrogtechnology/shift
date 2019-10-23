package cmd

import (
	"github.com/leapfrogtechnology/shift/cli/internals/env"
)

// AddEnv adds a new environment to the project.
func AddEnv() {
	env.Run()
}

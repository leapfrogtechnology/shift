package cmd

import (
	"github.com/leapfrogtechnology/shift/cli/internals/setup"
)

// Setup prompts user for required information for creating a project.
func Setup() {
	setup.Run()
}

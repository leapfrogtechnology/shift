package cmd

import "github.com/leapfrogtechnology/shift/cli/internals/destroy"

// Destroy delete existing infrastructure
func Destroy(environment string) {
	destroy.Run(environment)
}

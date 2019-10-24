package cmd

import "github.com/leapfrogtechnology/shift/cli/internals/destroy"

// Destroy deletes existing infrastructure
func Destroy(environment string) {
	destroy.Run(environment)
}

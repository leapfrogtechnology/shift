package infrastructure

import "github.com/leapfrogtechnology/shift/infrastructure/internals"

func Initialize(details string) (string, error) {
	infrastructureInfo, err := internals.Initialize(details)
	return infrastructureInfo, err
}

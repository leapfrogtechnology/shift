package infrastructure

import (
	"errors"
	"github.com/leapfrogtechnology/shift/core/utils/logger"
	"github.com/leapfrogtechnology/shift/infrastructure/internals"
	"github.com/leapfrogtechnology/shift/infrastructure/utils"
)

func Initialize(details string) (string, error) {
	if !utils.CommandExists("terraform") {
		logger.FailOnError(errors.New("terraform does not exist"), "Please install terraform on your device")
	}
	infrastructureInfo, err := internals.Initialize(details)
	if err != nil {
		logger.FailOnError(err, "Failed to Init Infrastructure")
	}
	return infrastructureInfo, err
}

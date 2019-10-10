package internals

import (
	"errors"
	"github.com/leapfrogtechnology/shift/core/utils/logger"
	"github.com/leapfrogtechnology/shift/infrastructure/infrastructure"
)

func Initialize(details string) (string, error) {
	infrastructureInfo, err := infrastrucuture.Initialize([]byte(details))
	if err != nil {
		logger.LogError(err, "Cannot Init Infrastructure")
		return "", errors.New("issue Initializing the infrastructure")
	} else {
		logger.LogOutput(infrastructureInfo)
		return infrastructureInfo, nil
	}
}

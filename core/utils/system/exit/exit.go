package exit

import (
	"os"

	"github.com/leapfrogtechnology/shift/core/utils/logger"
)

// Error exits the application by showing the given message.
func Error(err error, msg string) {
	logger.Error(err, msg)

	os.Exit(1)
}

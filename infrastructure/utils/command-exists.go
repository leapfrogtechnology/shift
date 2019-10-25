package utils

import (
	"os/exec"
)

// CommandExists checks if the command exists.
func CommandExists(command string) bool {
	_, err := exec.LookPath(command)

	return err == nil
}

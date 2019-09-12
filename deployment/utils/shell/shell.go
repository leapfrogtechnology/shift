package shell

import (
	"os"
	"os/exec"
)

// Execute runs the process with the supplied environment.
func Execute(command string) error {
	cmd := exec.Command("bash", "-c", command)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	return err
}

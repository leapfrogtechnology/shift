package shell

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// Execute runs the process with the supplied environment.
func Execute(command string) error {
	cmd := exec.Command("bash", "-c", command)

	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stderr = &stderr
	cmd.Stdout = &out

	cmd.Stdout = os.Stdout

	err := cmd.Run()

	if err != nil {
		return errors.New(fmt.Sprint(err) + ":- " + stderr.String())
	}

	return nil
}

package shell

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/leapfrogtechnology/shift/core/utils/escape"
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
		fmt.Println(stderr.String())

		errorMessage := escape.Strip(err.Error() + " :- " + stderr.String())

		return errors.New(errorMessage)
	}
	return nil
}

package cli

import (
	"fmt"

	"github.com/leapfrogtechnology/shift/cli/cmd"
)

func Initialize() {
	info := &cmd.Info{
		Name:        "Shift",
		Version:     "0.0.1",
		Description: "CLI for Shift",
	}

	err := cmd.Initialize(info)

	if err != nil {
		fmt.Println(err)
	}
}

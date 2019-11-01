package main

import (
	"fmt"

	"github.com/leapfrogtechnology/shift/cli/cmd"
)

// Initialize Cli
func main() {
	info := &cmd.Info{
		Name:        "Shift",
		Version:     "0.0.3",
		Description: "CLI for Shift",
	}

	err := cmd.Initialize(info)

	if err != nil {
		fmt.Println(err)
	}
}

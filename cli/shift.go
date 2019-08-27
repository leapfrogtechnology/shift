package main

import (
	"fmt"

	"github.com/lftechnology/shift/cli/cmd"
)

func main() {
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

package cmd

import (
	"os"

	"github.com/urfave/cli"
)

// Info defines the basic information required for the CLI.
type Info struct {
	Name        string
	Version     string
	Description string
}

// Initialize and bootstrap the CLI.
func Initialize(info *Info) error {
	app := cli.NewApp()

	app.Name = info.Name
	app.Version = info.Version
	app.Usage = info.Description

	app.Commands = []cli.Command{
		{
			Name:        "infrastructure",
			Description: "Initialize",
			Aliases:     nil,
			Usage:       "Initialize your Application",
			Action: func(c *cli.Context) {
				fmt.Println("Shift Shift shift")
			},
		},
		cli.Command{
			Name: "setup",
			Action: func(ctx *cli.Context) error {
				Setup()

				return nil
			},
		},
	}

	return app.Run(os.Args)
}

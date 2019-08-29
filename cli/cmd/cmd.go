package cmd

import (
	// "github.com/leapfrogtechnology/shift/infrastructure"
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
		cli.Command{
			Name:        "init",
			Description: "Initialize",
			Aliases:     nil,
			Usage:       "Initialize your Application",
			Subcommands: []cli.Command{
				{
					Name:        "frontend",
					Aliases:     nil,
					Usage:       "Initialize your frontend infrastructure",
					Description: "Use this to initialize your frontend Infrastructure",
					Action: func(c *cli.Context) {
						// infrastructure.InitializeFrontend()
					},
				},
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

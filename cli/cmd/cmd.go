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
		cli.Command{
			Name: "setup",
			Action: func(ctx *cli.Context) error {
				Setup()

				return nil
			},
		},
		cli.Command{
			Name: "deploy",
			Action: func(ctx *cli.Context) error {
				environment := ctx.Args().Get(0)
				Deploy(environment)

				return nil
			},
		},
		cli.Command{
			Name: "destroy",
			Action: func(ctx *cli.Context) error {
				environment := ctx.Args().Get(0)
				Destroy(environment)

				return nil
			},
		},
		cli.Command{
			Name: "add",
			Subcommands: []cli.Command{
				{
					Name: "env",
					Action: func(ctx *cli.Context) error {
						AddEnv()

						return nil
					},
				},
			},
		},
	}

	return app.Run(os.Args)
}

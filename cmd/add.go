package cmd

import "github.com/urfave/cli/v2"

var addCmd = &cli.Command{
	Name:  "add",
	Usage: "add framework support",
	Subcommands: []*cli.Command{
		addTypescriptCmd,
		addReactCmd,
	},
}

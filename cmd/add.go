package cmd

import "github.com/urfave/cli/v2"

var addCmd = &cli.Command{
	Name:    "add",
	Aliases: []string{"a"},
	Usage:   "add framework support",
	Subcommands: []*cli.Command{
		addTypescriptCmd,
		addReactCmd,
		addGoAirCmd,
		addGoMakefileCmd,
		addPriceUpdateConfigCmd,
		addViperTagCmd,
	},
}

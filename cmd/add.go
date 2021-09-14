package cmd

import "github.com/urfave/cli/v2"

var AddSet = &cli.Command{
	Name:    "add",
	Aliases: []string{"a"},
	Usage:   "add framework support",
	Subcommands: []*cli.Command{
		AddTypescriptCmd,
		AddReactCmd,
		AddGoAirCmd,
		AddGoMakefileCmd,
		AddPriceUpdateConfigCmd,
		AddViperTagCmd,
	},
}

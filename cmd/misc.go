package cmd

import "github.com/urfave/cli/v2"

var MiscSet = &cli.Command{
	Name:    "misc",
	Aliases: []string{"m"},
	Usage:   "un-categorized tools",
	Subcommands: []*cli.Command{
		MassRenameCmd,
	},
}

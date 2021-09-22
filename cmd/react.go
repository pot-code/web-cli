package cmd

import "github.com/urfave/cli/v2"

var ReactSet = &cli.Command{
	Name:    "react",
	Aliases: []string{"r"},
	Usage:   "react toolkit",
	Subcommands: []*cli.Command{
		ReactTemplateCmd,
		ReactComponentCmd,
	},
}
package cmd

import "github.com/urfave/cli/v2"

var generateCmd = &cli.Command{
	Name:    "generate",
	Aliases: []string{"g", "gen"},
	Usage:   "generate project files based on some templates",
	Subcommands: []*cli.Command{
		generateBECmd,
		genAPICmd,
	},
}

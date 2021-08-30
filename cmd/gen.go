package cmd

import "github.com/urfave/cli/v2"

var generateCmd = &cli.Command{
	Name:    "generate",
	Aliases: []string{"g", "gen"},
	Usage:   "generate module or project structure",
	Subcommands: []*cli.Command{
		generateBECmd,
		genAPICmd,
		generateFECmd,
	},
}

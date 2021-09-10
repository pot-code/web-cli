package cmd

import "github.com/urfave/cli/v2"

var generateCmd = &cli.Command{
	Name:    "generate",
	Aliases: []string{"g", "gen"},
	Usage:   "generate modules or project",
	Subcommands: []*cli.Command{
		generateBECmd,
		genAPICmd,
		generateFECmd,
		genFlagsCmd,
	},
}

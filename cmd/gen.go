package cmd

import "github.com/urfave/cli/v2"

var GenerationSet = &cli.Command{
	Name:    "generate",
	Aliases: []string{"g", "gen"},
	Usage:   "generate modules or project structure",
	Subcommands: []*cli.Command{
		GenBECmd,
		GenAPICmd,
		GenFECmd,
		GenFlagsCmd,
	},
}

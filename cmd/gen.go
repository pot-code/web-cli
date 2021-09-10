package cmd

import "github.com/urfave/cli/v2"

var GenerateCmd = &cli.Command{
	Name:    "generate",
	Aliases: []string{"g", "gen"},
	Usage:   "generate modules or project",
	Subcommands: []*cli.Command{
		GenBeCmd,
		GenAPICmd,
		GenerateFECmd,
		GenFlagsCmd,
	},
}

package cmd

import "github.com/urfave/cli/v2"

var Generate = &cli.Command{
	Name:    "generate",
	Aliases: []string{"gen"},
	Usage:   "generate project files based on some templates",
	Subcommands: []*cli.Command{
		genBE,
		genAPI,
	},
}

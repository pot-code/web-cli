package admingo

import "github.com/urfave/cli/v2"

var CommandSet = &cli.Command{
	Name:    "admingo",
	Aliases: []string{"ag"},
	Usage:   "fast admin go 套件",
	Subcommands: []*cli.Command{
		CreateModuleCmd,
	},
}

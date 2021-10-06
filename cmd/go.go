package cmd

import "github.com/urfave/cli/v2"

var GolangSet = &cli.Command{
	Name:    "go",
	Aliases: []string{"g"},
	Usage:   "golang toolkit",
	Subcommands: []*cli.Command{
		GoWebCmd,
		GoServiceCmd,
		GoFlagsCmd,
		GoAirCmd,
		GoMakefileCmd,
		GoViperTagCmd,
		GoMigrateCmd,
	},
}

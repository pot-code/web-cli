package cmd

import "github.com/urfave/cli/v2"

var VueSet = &cli.Command{
	Name:    "vue",
	Aliases: []string{"v"},
	Usage:   "vue kit",
	Subcommands: []*cli.Command{
		VueUseStoreCmd,
		VueComponentCmd,
		VueDirectiveCmd,
	},
}

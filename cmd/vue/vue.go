package vue

import "github.com/urfave/cli/v2"

var CommandSet = &cli.Command{
	Name:    "vue",
	Aliases: []string{"v"},
	Usage:   "vue kit",
	Subcommands: []*cli.Command{
		VueUseStoreCmd,
		VueComponentCmd,
		VueDirectiveCmd,
	},
}

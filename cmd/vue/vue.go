package vue

import "github.com/urfave/cli/v2"

var CommandSet = &cli.Command{
	Name:    "vue",
	Aliases: []string{"v"},
	Usage:   "vue 工具套件",
	Subcommands: []*cli.Command{
		VueUseStoreCmd,
		VueComponentCmd,
		VueDirectiveCmd,
	},
}

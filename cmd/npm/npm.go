package npm

import "github.com/urfave/cli/v2"

var CommandSet = &cli.Command{
	Name:    "npm",
	Aliases: []string{"n"},
	Usage:   "npm 包管理工具",
	Subcommands: []*cli.Command{
		PinPackageCmd,
	},
}

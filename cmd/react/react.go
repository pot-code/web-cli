package react

import "github.com/urfave/cli/v2"

var CommandSet = &cli.Command{
	Name:    "react",
	Aliases: []string{"r"},
	Usage:   "react kit",
	Subcommands: []*cli.Command{
		ReactComponentCmd,
		ReactHookCmd,
		ReactZustandCmd,
		ReactContextCmd,
	},
}

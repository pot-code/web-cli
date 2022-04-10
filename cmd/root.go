package cmd

import (
	"github.com/urfave/cli/v2"
)

var AppVersion string

var RootCmd = &cli.App{
	Name:  "web-cli",
	Usage: "A toolset to improve your development experience",
	Commands: []*cli.Command{
		GolangSet,
		ReactSet,
		MiscSet,
	},
	Version: AppVersion,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:        "debug",
			Aliases:     []string{"D"},
			Usage:       "enable debugging",
			Value:       false,
			DefaultText: "false",
		},
	},
	Before: func(c *cli.Context) error {
		// if c.Bool("debug") {
		// 	log.SetLevel(log.DebugLevel)
		// }
		return nil
	},
}

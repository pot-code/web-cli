package cmd

import (
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

var AppVersion string

var RootCmd = &cli.App{
	Name:  "web-cli",
	Usage: "A toolkit to for web development",
	Commands: []*cli.Command{
		ReactSet,
		VueSet,
		MiscSet,
	},
	Version: AppVersion,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:        "verbose",
			Aliases:     []string{"vv"},
			Usage:       "verbose mode",
			Value:       false,
			DefaultText: "false",
		},
	},
	Before: func(c *cli.Context) error {
		if c.Bool("verbose") {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		}
		return nil
	},
}

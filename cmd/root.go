package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var AppVersion string

var RootCmd = &cli.App{
	Name:  "web-cli",
	Usage: "A toolkit to for web development",
	Commands: []*cli.Command{
		ReactSet,
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
			log.SetLevel(log.DebugLevel)
		}
		return nil
	},
}

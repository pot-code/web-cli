package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var AppVersion string

var RootCmd = &cli.App{
	Name:  "web-cli",
	Usage: "web utils",
	Commands: []*cli.Command{
		Generate,
	},
	Version: AppVersion,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:        "debug",
			Aliases:     []string{"d"},
			Usage:       "enable debug",
			Value:       false,
			DefaultText: "false",
		},
	},
	Before: func(c *cli.Context) error {
		if c.Bool("debug") {
			log.SetLevel(log.DebugLevel)
		}
		return nil
	},
}

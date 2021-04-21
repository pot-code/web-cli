package cmd

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/pot-code/web-cli/add"
	"github.com/pot-code/web-cli/core"
	"github.com/pot-code/web-cli/util"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const cmdReactName = "React"

type addReactConfig struct {
	Hook bool
	Name string
}

var addReactCmd = &cli.Command{
	Name:      cmdReactName,
	Usage:     "add React components",
	Aliases:   []string{"r"},
	ArgsUsage: "NAME",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "hook",
			Aliases: []string{"H"},
			Usage:   "add hook",
		},
	},
	Action: func(c *cli.Context) error {
		config, err := getAddReactConfig(c)
		if err != nil {
			if _, ok := err.(*util.CommandError); ok {
				cli.ShowCommandHelp(c, cmdReactName)
			}
			return err
		}

		var cmd core.Executor
		if config.Hook {
			cmd = add.NewAddReactHook(config.Name)
		} else {
			cmd = add.NewAddReactComponent(config.Name)
		}
		return cmd.Run()
	},
}

func getAddReactConfig(c *cli.Context) (*addReactConfig, error) {
	name := c.Args().Get(0)
	if name == "" {
		return nil, util.NewCommandError(cmdReactName, fmt.Errorf("NAME must be specified"))
	}
	if err := util.ValidateVarName(name); err != nil {
		return nil, util.NewCommandError(cmdBackendName, errors.Wrap(err, "invalid NAME"))
	}
	name = strings.ReplaceAll(name, "-", "_")

	hook := c.Bool("hook")
	if hook {
		name = strcase.ToLowerCamel(name)
	} else {
		name = strcase.ToCamel(name)
	}
	log.Debug("transformed component name: ", name)

	config := &addReactConfig{Hook: hook, Name: name}
	return config, nil
}

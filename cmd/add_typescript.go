package cmd

import (
	"fmt"
	"strings"

	"github.com/pot-code/web-cli/add"
	"github.com/pot-code/web-cli/core"
	"github.com/pot-code/web-cli/util"
	"github.com/urfave/cli/v2"
)

const cmdTypescriptName = "typescript"

type addTypescriptConfig struct {
	Target string
}

var addTypescriptCmd = &cli.Command{
	Name:    cmdTypescriptName,
	Usage:   "add typescript support",
	Aliases: []string{"ts"},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "target",
			Aliases: []string{"T"},
			Usage:   "project target (node/react)",
			Value:   "node",
		},
	},
	Action: func(c *cli.Context) error {
		config, err := getAddTypescriptConfig(c)
		if err != nil {
			if _, ok := err.(*util.CommandError); ok {
				cli.ShowCommandHelp(c, cmdTypescriptName)
			}
			return err
		}

		var cmd core.Generator
		if config.Target == "node" {
			cmd = add.NewAddTypescriptToNode()
		} else {
			cmd = add.NewAddTypescriptToReact()
		}

		err = cmd.Run()
		if err != nil {
			cmd.Cleanup()
		}
		return err
	},
}

func getAddTypescriptConfig(c *cli.Context) (*addTypescriptConfig, error) {
	config := &addTypescriptConfig{}
	target := strings.ToLower(c.String("target"))
	if target == "" {
		return nil, util.NewCommandError(cmdTypescriptName, fmt.Errorf("target is empty"))
	} else if target == "node" || target == "react" {
		config.Target = target
	} else {
		return nil, util.NewCommandError(cmdTypescriptName, fmt.Errorf("unsupported target '%s'", target))
	}
	return config, nil
}

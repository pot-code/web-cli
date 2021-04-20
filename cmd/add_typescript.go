package cmd

import (
	"fmt"
	"strings"

	"github.com/pot-code/web-cli/add"
	"github.com/pot-code/web-cli/util"
	"github.com/urfave/cli/v2"
)

const CmdTypescriptName = "typescript"

type AddTypescriptConfig struct {
	Target string
}

var addTypescriptCmd = &cli.Command{
	Name:    CmdTypescriptName,
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
				cli.ShowCommandHelp(c, CmdTypescriptName)
			}
			return err
		}

		if config.Target == "node" {
			cmd := add.NewAddTypescriptToNode()
			return cmd.Run()
		}
		cmd := add.NewAddTypescriptToReact()
		return cmd.Run()
	},
}

func getAddTypescriptConfig(c *cli.Context) (*AddTypescriptConfig, error) {
	config := &AddTypescriptConfig{}
	target := strings.ToLower(c.String("target"))
	if target == "" {
		return nil, util.NewCommandError(CmdTypescriptName, fmt.Errorf("target is empty"))
	} else if target == "node" || target == "react" {
		config.Target = target
	} else {
		return nil, util.NewCommandError(CmdTypescriptName, fmt.Errorf("unsupported target '%s'", target))
	}
	return config, nil
}

package cmd

import (
	"github.com/pot-code/web-cli/pkg/core"
	"github.com/pot-code/web-cli/pkg/util"
	"github.com/urfave/cli/v2"
)

type genFEConfig struct {
	GenType     string `flag:"type" validate:"required,oneof=react next"`     // generation type
	ProjectName string `arg:"0" alias:"project_name" validate:"required,var"` // project name
}

var generateFECmd = &cli.Command{
	Name:      "frontend",
	Aliases:   []string{"fe"},
	Usage:     "generate frontends",
	ArgsUsage: "project_name",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "type",
			Aliases: []string{"t"},
			Usage:   "framework type (react)",
			Value:   "react",
		},
	},
	Action: func(c *cli.Context) error {
		config := new(genFEConfig)
		err := util.ParseConfig(c, config)
		if err != nil {
			if _, ok := err.(*util.CommandError); ok {
				cli.ShowCommandHelp(c, c.Command.Name)
			}
			return err
		}

		if config.GenType == "react" {
			gen := generateReact(config)
			if err := gen.Run(); err != nil {
				return err
			}
		}

		if config.GenType == "next" {
			gen := generateNext(config)
			if err := gen.Run(); err != nil {
				return err
			}
		}
		return nil
	},
}

func generateReact(config *genFEConfig) core.Runner {
	return util.NewTaskComposer("").AddCommand(&core.Command{
		Bin:  "git",
		Args: []string{"clone", "https://github.com/pot-code/react-boilerplate.git", config.ProjectName},
	})
}

func generateNext(config *genFEConfig) core.Runner {
	return util.NewTaskComposer("").AddCommand(&core.Command{
		Bin:  "yarn",
		Args: []string{"create", "next-app", "--ts", config.ProjectName},
	})
}

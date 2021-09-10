package cmd

import (
	"github.com/pot-code/web-cli/pkg/core"
	"github.com/pot-code/web-cli/pkg/util"
	"github.com/urfave/cli/v2"
)

type GenFEConfig struct {
	GenType     string `flag:"type" alias:"t" usage:"framework type" validate:"required,oneof=react next"` // generation type
	ProjectName string `arg:"0" alias:"project_name" validate:"required,var"`                              // project name
}

var GenerateFECmd = core.NewCliLeafCommand("frontend", "generate frontends",
	&GenFEConfig{
		GenType: "react",
	},
	core.WithArgUsage("project_name"),
	core.WithAlias([]string{"fe"}),
).AddService(
	new(GenerateReactFeService),
	new(GenerateNextJsFeService),
).ExportCommand()

type GenerateReactFeService struct{}

var _ core.CommandService = &GenerateReactFeService{}

func (grf *GenerateReactFeService) Cond(c *cli.Context) bool {
	return c.String("type") == "react"
}

func (grf *GenerateReactFeService) Handle(c *cli.Context, cfg interface{}) error {
	config := cfg.(*GenFEConfig)

	return util.NewTaskComposer("").AddCommand(
		&core.ShellCommand{
			Bin:  "git",
			Args: []string{"clone", "https://github.com/pot-code/react-boilerplate.git", config.ProjectName},
		},
	).Run()
}

type GenerateNextJsFeService struct{}

var _ core.CommandService = &GenerateNextJsFeService{}

func (gnf *GenerateNextJsFeService) Cond(c *cli.Context) bool {
	return c.String("type") == "next"
}

func (gnf *GenerateNextJsFeService) Handle(c *cli.Context, cfg interface{}) error {
	config := cfg.(*GenFEConfig)

	return util.NewTaskComposer("").AddCommand(
		&core.ShellCommand{
			Bin:  "yarn",
			Args: []string{"create", "next-app", "--ts", config.ProjectName},
		},
	).Run()
}

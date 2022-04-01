package cmd

import (
	"github.com/pot-code/web-cli/pkg/commands"
	"github.com/pot-code/web-cli/pkg/core"
	"github.com/pot-code/web-cli/pkg/util"
	"github.com/urfave/cli/v2"
)

type ReactInitConfig struct {
	GenType     string `flag:"type" alias:"t" usage:"project type" validate:"required,oneof=vanilla next"`
	ProjectName string `arg:"0" alias:"project_name" validate:"required,var"`
}

var ReactInitCmd = core.NewCliLeafCommand("init", "create react project",
	&ReactInitConfig{
		GenType: "vanilla",
	},
	core.WithArgUsage("project_name"),
).AddFeature(
	new(VanillaTemplate),
	new(NextJsTemplate),
).ExportCommand()

type VanillaTemplate struct{}

var _ core.CommandFeature = &VanillaTemplate{}

func (grf *VanillaTemplate) Cond(c *cli.Context) bool {
	return c.String("type") == "vanilla"
}

func (grf *VanillaTemplate) Handle(c *cli.Context, cfg interface{}) error {
	config := cfg.(*ReactInitConfig)

	return util.NewTaskComposer("").AddCommand(
		commands.GitClone("https://github.com/pot-code/react-boilerplate.git", config.ProjectName),
		commands.GitDeleteHistory(config.ProjectName),
	).Run()
}

type NextJsTemplate struct{}

var _ core.CommandFeature = &NextJsTemplate{}

func (gnf *NextJsTemplate) Cond(c *cli.Context) bool {
	return c.String("type") == "next"
}

func (gnf *NextJsTemplate) Handle(c *cli.Context, cfg interface{}) error {
	config := cfg.(*ReactInitConfig)

	return util.NewTaskComposer("").AddCommand(
		commands.YarnCreate("next-app", "--ts", config.ProjectName),
	).Run()
}

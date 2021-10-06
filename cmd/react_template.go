package cmd

import (
	"github.com/pot-code/web-cli/pkg/commands"
	"github.com/pot-code/web-cli/pkg/core"
	"github.com/pot-code/web-cli/pkg/util"
	"github.com/urfave/cli/v2"
)

type ReactTemplateConfig struct {
	GenType     string `flag:"type" alias:"t" usage:"framework type" validate:"required,oneof=react next"` // generation type
	ProjectName string `arg:"0" alias:"project_name" validate:"required,var"`                              // project name
}

var ReactTemplateCmd = core.NewCliLeafCommand("template", "choose react template",
	&ReactTemplateConfig{
		GenType: "react",
	},
	core.WithArgUsage("project_name"),
	core.WithAlias([]string{"t"}),
).AddFeature(
	new(GenReactTemplate),
	new(GenNextJsTemplate),
).ExportCommand()

type GenReactTemplate struct{}

var _ core.CommandFeature = &GenReactTemplate{}

func (grf *GenReactTemplate) Cond(c *cli.Context) bool {
	return c.String("type") == "react"
}

func (grf *GenReactTemplate) Handle(c *cli.Context, cfg interface{}) error {
	config := cfg.(*ReactTemplateConfig)

	return util.NewTaskComposer("").AddCommand(
		commands.GitClone("https://github.com/pot-code/react-boilerplate.git", config.ProjectName),
	).Run()
}

type GenNextJsTemplate struct{}

var _ core.CommandFeature = &GenNextJsTemplate{}

func (gnf *GenNextJsTemplate) Cond(c *cli.Context) bool {
	return c.String("type") == "next"
}

func (gnf *GenNextJsTemplate) Handle(c *cli.Context, cfg interface{}) error {
	config := cfg.(*ReactTemplateConfig)

	return util.NewTaskComposer("").AddCommand(
		commands.YarnCreate("next-app", "--ts", config.ProjectName),
	).Run()
}

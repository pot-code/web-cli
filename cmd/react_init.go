package cmd

import (
	"github.com/pot-code/web-cli/internal/command"
	"github.com/pot-code/web-cli/internal/shell"
	"github.com/pot-code/web-cli/internal/task"
	"github.com/urfave/cli/v2"
)

type ReactInitConfig struct {
	GenType     string `flag:"type" alias:"t" usage:"project type" validate:"oneof=vanilla next"`
	ProjectName string `arg:"0" alias:"project_name" validate:"required,var"`
}

var ReactInitCmd = command.NewCliCommand("init", "create react project",
	&ReactInitConfig{
		GenType: "vanilla",
	},
	command.WithArgUsage("project_name"),
).AddHandlers(
	new(VanillaTemplate),
	new(NextJsTemplate),
).BuildCommand()

type VanillaTemplate struct{}

var _ command.CommandHandler = &VanillaTemplate{}

func (grf *VanillaTemplate) Cond(c *cli.Context) bool {
	return c.String("type") == "vanilla"
}

func (grf *VanillaTemplate) Handle(c *cli.Context, cfg interface{}) error {
	config := cfg.(*ReactInitConfig)

	return task.NewSequentialExecutor(
		[]task.Task{
			shell.GitClone("https://github.com/pot-code/react-boilerplate.git", config.ProjectName),
			shell.GitDeleteHistory(config.ProjectName),
		},
	).Run()
}

type NextJsTemplate struct{}

var _ command.CommandHandler = &NextJsTemplate{}

func (gnf *NextJsTemplate) Cond(c *cli.Context) bool {
	return c.String("type") == "next"
}

func (gnf *NextJsTemplate) Handle(c *cli.Context, cfg interface{}) error {
	config := cfg.(*ReactInitConfig)

	return shell.YarnCreate("next-app", "--ts", config.ProjectName).Run()
}

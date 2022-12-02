package cmd

import (
	"github.com/iancoleman/strcase"
	"github.com/pot-code/web-cli/internal/command"
	"github.com/pot-code/web-cli/internal/env"
	"github.com/pot-code/web-cli/internal/task"
	"github.com/urfave/cli/v2"
)

const TypescriptSuffix = ".ts"

type ReactHookConfig struct {
	Name   string `arg:"0" alias:"HOOK_NAME" validate:"required,nature"`
	OutDir string `flag:"output" alias:"o" usage:"destination directory"`
}

var ReactHookCmd = command.NewCliCommand("hook", "add react hook",
	new(ReactHookConfig),
	command.WithArgUsage("HOOK_NAME"),
	command.WithAlias([]string{"k"}),
).AddHandlers(
	AddReactHook,
).BuildCommand()

var AddReactHook = command.InlineHandler(func(c *cli.Context, cfg interface{}) error {
	config := cfg.(*ReactHookConfig)
	name := strcase.ToLowerCamel(config.Name)

	return task.NewSequentialExecutor().AddTask(
		task.NewTemplateRenderTask(
			name,
			TypescriptSuffix,
			config.OutDir,
			task.NewLocalTemplateProvider(env.GetAbsoluteTemplatePath("react_hook.tmpl")),
			false,
			map[string]string{
				"name": name,
			},
		),
	).Run()
})

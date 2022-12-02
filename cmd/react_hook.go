package cmd

import (
	"github.com/iancoleman/strcase"
	"github.com/pot-code/web-cli/internal/command"
	"github.com/pot-code/web-cli/internal/task"
	"github.com/pot-code/web-cli/internal/template"
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

	return task.NewGenerateFileFromTemplateTask(
		name,
		TypescriptSuffix,
		config.OutDir,
		false,
		name,
		template.NewLocalTemplateProvider(GetAbsoluteTemplatePath("react_hook.tmpl")),
		map[string]string{
			"name": name,
		},
	).Run()
})

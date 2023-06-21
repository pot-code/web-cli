package cmd

import (
	"bytes"
	"github.com/iancoleman/strcase"
	"github.com/pot-code/web-cli/pkg/command"
	"github.com/pot-code/web-cli/pkg/file"
	"github.com/pot-code/web-cli/pkg/provider"
	"github.com/pot-code/web-cli/pkg/task"
	"github.com/urfave/cli/v2"
)

type ReactContextConfig struct {
	Name   string `arg:"0" alias:"CONTEXT_NAME" validate:"required,var"`
	OutDir string `flag:"output" alias:"o" usage:"destination directory"`
}

var ReactContextCmd = command.NewCliCommand("context", "add custom context",
	new(ReactContextConfig),
	command.WithArgUsage("CONTEXT_NAME"),
	command.WithAlias([]string{"ctx"}),
).AddHandlers(
	AddContextStore,
).BuildCommand()

var AddContextStore = command.InlineHandler(func(c *cli.Context, cfg interface{}) error {
	rzc := cfg.(*ReactContextConfig)
	varName := strcase.ToCamel(rzc.Name)
	filename := strcase.ToKebab(varName)

	b := new(bytes.Buffer)
	tasks := []task.Task{
		task.NewSequentialScheduler().
			AddTask(task.NewReadFromProviderTask(provider.NewEmbedFileProvider("templates/react/react_context.gotmpl"), b)).
			AddTask(task.NewTemplateRenderTask("react_context", map[string]string{"name": varName}, b, b)).
			AddTask(task.NewWriteFileToDiskTask(filename, file.ReactComponentSuffix, rzc.OutDir, false, b)),
	}

	s := task.NewParallelScheduler()
	for _, t := range tasks {
		s.AddTask(t)
	}
	return s.Run()
})

package react

import (
	"bytes"

	"github.com/iancoleman/strcase"
	"github.com/pot-code/web-cli/internal/command"
	"github.com/pot-code/web-cli/internal/provider"
	"github.com/pot-code/web-cli/internal/task"
	"github.com/urfave/cli/v2"
)

type ReactContextConfig struct {
	Name   string `arg:"0" alias:"CONTEXT_NAME" validate:"required,var"`
	OutDir string `flag:"output" alias:"o" usage:"destination directory"`
}

var ReactContextCmd = command.NewBuilder("context", "add custom context",
	new(ReactContextConfig),
	command.WithArgUsage("CONTEXT_NAME"),
	command.WithAlias([]string{"ctx"}),
).AddHandlers(
	AddContextStore,
).Build()

var AddContextStore = command.InlineHandler(func(c *cli.Context, cfg interface{}) error {
	rzc := cfg.(*ReactContextConfig)
	filename := strcase.ToKebab(rzc.Name)
	ctxName := strcase.ToCamel(rzc.Name)

	b := new(bytes.Buffer)
	tasks := []task.Task{
		task.NewSequentialScheduler().
			AddTask(task.NewReadFromProviderTask(provider.NewEmbedFileProvider("templates/react/react_context.gotmpl"), b)).
			AddTask(task.NewTemplateRenderTask("react_context", map[string]string{"name": ctxName}, b, b)).
			AddTask(task.NewWriteFileToDiskTask(filename, ".tsx", rzc.OutDir, false, b)),
	}

	s := task.NewParallelScheduler()
	for _, t := range tasks {
		s.AddTask(t)
	}
	return s.Run()
})

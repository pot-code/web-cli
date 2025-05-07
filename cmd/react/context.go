package react

import (
	"bytes"
	"fmt"

	"github.com/iancoleman/strcase"
	"github.com/pot-code/web-cli/internal/command"
	"github.com/pot-code/web-cli/internal/provider"
	"github.com/pot-code/web-cli/internal/task"
	"github.com/urfave/cli/v2"
)

type ReactContextConfig struct {
	Name   string `arg:"0" alias:"CONTEXT_NAME" validate:"required,identifier"`
	OutDir string `flag:"output" alias:"o" usage:"destination directory"`
}

var ReactContextCmd = command.NewCommand("context", "add custom context",
	new(ReactContextConfig),
	command.WithAlias([]string{"ctx"}),
).AddHandler(
	AddContextStore,
).Create()

var AddContextStore = command.InlineHandler[*ReactContextConfig](func(c *cli.Context, config *ReactContextConfig) error {
	contextFileName := strcase.ToKebab(fmt.Sprintf("%sContext", config.Name))
	contextName := strcase.ToCamel(fmt.Sprintf("%sContext", config.Name))
	hookFileName := strcase.ToKebab(fmt.Sprintf("use%s", contextName))

	b := new(bytes.Buffer)
	b1 := new(bytes.Buffer)
	tasks := []task.Task{
		task.NewSequentialScheduler().
			AddTask(task.NewReadFromProviderTask(provider.NewEmbedFileProvider("templates/react/react_context.go.tmpl"), b)).
			AddTask(task.NewTemplateRenderTask("react_context", map[string]string{"name": contextName}, b, b)).
			AddTask(task.NewWriteFileToDiskTask(contextFileName, ".tsx", config.OutDir, false, b)),
		task.NewSequentialScheduler().
			AddTask(task.NewReadFromProviderTask(provider.NewEmbedFileProvider("templates/react/react_context_hook.go.tmpl"), b1)).
			AddTask(task.NewTemplateRenderTask("react_context_hook", map[string]string{"name": contextName, "file": contextFileName}, b1, b1)).
			AddTask(task.NewWriteFileToDiskTask(hookFileName, ".ts", config.OutDir, false, b1)),
	}

	s := task.NewParallelScheduler()
	for _, t := range tasks {
		s.AddTask(t)
	}
	return s.Run()
})

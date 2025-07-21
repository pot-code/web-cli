package react

import (
	"bytes"
	"fmt"

	"github.com/iancoleman/strcase"
	"github.com/pot-code/web-cli/pkg/command"
	"github.com/pot-code/web-cli/pkg/provider"
	"github.com/pot-code/web-cli/pkg/task"
	"github.com/urfave/cli/v2"
)

type ReactContextConfig struct {
	Name   string `arg:"0" alias:"CONTEXT_NAME" validate:"required,min=1,max=32,identifier"`
	OutDir string `flag:"output" alias:"o" usage:"destination directory"`
}

var ReactContextCmd = command.NewCommand("context", "add custom context",
	new(ReactContextConfig),
	command.WithAlias([]string{"ctx"}),
).AddHandler(
	AddContextStore,
).Create()

var AddContextStore = command.InlineHandler[*ReactContextConfig](func(c *cli.Context, config *ReactContextConfig) error {
	contextName := strcase.ToCamel(fmt.Sprintf("%sContext", config.Name))
	contextFileName := strcase.ToKebab(fmt.Sprintf("%sContext", config.Name))
	hookFileName := strcase.ToKebab(fmt.Sprintf("use%s", contextName))

	var buffers []*bytes.Buffer
	for range 2 {
		buffers = append(buffers, new(bytes.Buffer))
	}
	tasks := []task.Task{
		task.NewSequentialScheduler().
			AddTask(task.NewReadFromProviderTask(provider.NewEmbedFileProvider("templates/react/react_context.go.tmpl"), buffers[0])).
			AddTask(task.NewTemplateRenderTask("react_context", map[string]string{"name": contextName}, buffers[0], buffers[0])).
			AddTask(task.NewWriteFileToDiskTask(contextFileName, ".tsx", buffers[0], task.WithFolder(config.OutDir))),
		task.NewSequentialScheduler().
			AddTask(task.NewReadFromProviderTask(provider.NewEmbedFileProvider("templates/react/react_context_hook.go.tmpl"), buffers[1])).
			AddTask(task.NewTemplateRenderTask("react_context_hook", map[string]string{"name": contextName, "file": contextFileName}, buffers[1], buffers[1])).
			AddTask(task.NewWriteFileToDiskTask(hookFileName, ".ts", buffers[1], task.WithFolder(config.OutDir))),
	}

	s := task.NewParallelScheduler()
	for _, t := range tasks {
		s.AddTask(t)
	}
	return s.Run()
})

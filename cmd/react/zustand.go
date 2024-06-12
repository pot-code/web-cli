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

type ReactZustandConfig struct {
	Name   string `arg:"0" alias:"STORE_NAME" validate:"required,var"`
	OutDir string `flag:"output" alias:"o" usage:"destination directory"`
}

var ReactZustandCmd = command.NewBuilder("zustand", "add zustand store",
	new(ReactZustandConfig),
	command.WithAlias([]string{"z"}),
).AddHandlers(
	AddZustandStore,
).Build()

var AddZustandStore = command.InlineHandler[*ReactZustandConfig](func(c *cli.Context, config *ReactZustandConfig) error {
	varName := fmt.Sprintf("use%sStore", strcase.ToCamel(config.Name))
	filename := strcase.ToKebab(varName)

	b := new(bytes.Buffer)
	tasks := []task.Task{
		task.NewSequentialScheduler().
			AddTask(task.NewReadFromProviderTask(provider.NewEmbedFileProvider("templates/react/react_zustand.gotmpl"), b)).
			AddTask(task.NewTemplateRenderTask("react_zustand", map[string]string{"name": filename}, b, b)).
			AddTask(task.NewWriteFileToDiskTask(filename, ".ts", config.OutDir, false, b)),
	}

	s := task.NewParallelScheduler()
	for _, t := range tasks {
		s.AddTask(t)
	}
	return s.Run()
})

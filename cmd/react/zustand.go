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

type ReactZustandConfig struct {
	Name   string `arg:"0" alias:"STORE_NAME" validate:"required,identifier"`
	OutDir string `flag:"output" alias:"o" usage:"destination directory"`
}

var ReactZustandCmd = command.NewCommand("zustand", "add zustand store",
	new(ReactZustandConfig),
	command.WithAlias([]string{"z"}),
).AddHandler(
	AddZustandStore,
).Create()

var AddZustandStore = command.InlineHandler[*ReactZustandConfig](func(c *cli.Context, config *ReactZustandConfig) error {
	varName := fmt.Sprintf("use%sStore", strcase.ToCamel(config.Name))
	filename := strcase.ToKebab(varName)

	b := new(bytes.Buffer)
	tasks := []task.Task{
		task.NewSequentialScheduler().
			AddTask(task.NewReadFromProviderTask(provider.NewEmbedFileProvider("templates/react/react_zustand.go.tmpl"), b)).
			AddTask(task.NewTemplateRenderTask("react_zustand", map[string]string{"name": filename}, b, b)).
			AddTask(task.NewWriteFileToDiskTask(filename, ".ts", b, task.WithFolder(config.OutDir))),
	}

	s := task.NewParallelScheduler()
	for _, t := range tasks {
		s.AddTask(t)
	}
	return s.Run()
})

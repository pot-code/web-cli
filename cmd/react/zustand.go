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

var ReactZustandCmd = command.NewCommandBuilder("zustand", "add zustand store",
	new(ReactZustandConfig),
	command.WithArgUsage("STORE_NAME"),
	command.WithAlias([]string{"z"}),
).AddHandlers(
	AddZustandStore,
).Build()

var AddZustandStore = command.InlineHandler(func(c *cli.Context, cfg interface{}) error {
	rzc := cfg.(*ReactZustandConfig)
	varName := strcase.ToCamel(rzc.Name)
	filename := strcase.ToLowerCamel(fmt.Sprintf("use%sStore", varName))

	b := new(bytes.Buffer)
	tasks := []task.Task{
		task.NewSequentialScheduler().
			AddTask(task.NewReadFromProviderTask(provider.NewEmbedFileProvider("templates/react/react_zustand.gotmpl"), b)).
			AddTask(task.NewTemplateRenderTask("react_zustand", map[string]string{"name": varName}, b, b)).
			AddTask(task.NewWriteFileToDiskTask(filename, ".ts", rzc.OutDir, false, b)),
	}

	s := task.NewParallelScheduler()
	for _, t := range tasks {
		s.AddTask(t)
	}
	return s.Run()
})

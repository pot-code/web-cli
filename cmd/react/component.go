package react

import (
	"bytes"
	"path"

	"github.com/iancoleman/strcase"
	"github.com/pot-code/web-cli/internal/command"
	"github.com/pot-code/web-cli/internal/provider"
	"github.com/pot-code/web-cli/internal/task"
	"github.com/urfave/cli/v2"
)

type ReactComponentConfig struct {
	AddStory bool   `flag:"add-storybook" alias:"s" usage:"add storybook"`
	Isolated bool   `flag:"isolated" alias:"i" usage:"generate files in a folder"`
	OutDir   string `flag:"output" alias:"o" usage:"destination directory"`
	Name     string `arg:"0" alias:"COMPONENT_NAME" validate:"required,min=1,max=32,var"`
}

var ReactComponentCmd = command.NewBuilder("component", "add react component",
	new(ReactComponentConfig),
	command.WithAlias([]string{"c"}),
).AddHandlers(
	new(AddReactComponent),
).Build()

type AddReactComponent struct {
	tasks []task.Task
}

func (arc *AddReactComponent) Handle(c *cli.Context, config *ReactComponentConfig) error {
	filename := strcase.ToKebab(config.Name)
	varName := strcase.ToCamel(config.Name)
	outDir := config.OutDir

	if config.Isolated {
		outDir = path.Join(outDir, filename)
		filename = "index"
	}
	arc.addComponent(varName, filename, outDir)

	if config.AddStory {
		arc.addStory(varName, filename, outDir)
	}

	e := task.NewParallelScheduler()
	for _, t := range arc.tasks {
		e.AddTask(t)
	}
	if err := e.Run(); err != nil {
		return err
	}
	return nil
}

func (arc *AddReactComponent) addComponent(varName string, filename string, outDir string) {
	b := new(bytes.Buffer)
	arc.tasks = append(arc.tasks,
		task.NewSequentialScheduler().
			AddTask(task.NewReadFromProviderTask(provider.NewEmbedFileProvider("templates/react/react_component.gotmpl"), b)).
			AddTask(task.NewTemplateRenderTask("react_component", map[string]string{"name": varName}, b, b)).
			AddTask(task.NewWriteFileToDiskTask(filename, ".tsx", outDir, false, b)))
}

func (arc *AddReactComponent) addStory(componentName string, filename string, outDir string) {
	b := new(bytes.Buffer)
	arc.tasks = append(arc.tasks,
		task.NewSequentialScheduler().
			AddTask(task.NewReadFromProviderTask(provider.NewEmbedFileProvider("templates/react/react_storybook.gotmpl"), b)).
			AddTask(task.NewTemplateRenderTask(
				"react_storybook",
				map[string]string{"name": componentName, "file": filename},
				b, b)).
			AddTask(task.NewWriteFileToDiskTask(filename, ".stories.tsx", outDir, false, b)))
}

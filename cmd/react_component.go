package cmd

import (
	"bytes"
	"path"

	"github.com/iancoleman/strcase"
	"github.com/pot-code/web-cli/pkg/command"
	"github.com/pot-code/web-cli/pkg/file"
	"github.com/pot-code/web-cli/pkg/provider"
	"github.com/pot-code/web-cli/pkg/task"
	"github.com/urfave/cli/v2"
)

type ReactComponentConfig struct {
	AddStory  bool   `flag:"add-storybook" alias:"s" usage:"add storybook"`
	AddFolder bool   `flag:"add-folder" alias:"f" usage:"generate files in a folder with the name as the component"`
	OutDir    string `flag:"output" alias:"o" usage:"destination directory"`
	Name      string `arg:"0" alias:"COMPONENT_NAME" validate:"required,var"`
}

var ReactComponentCmd = command.NewCliCommand("component", "add react component",
	new(ReactComponentConfig),
	command.WithArgUsage("COMPONENT_NAME"),
	command.WithAlias([]string{"c"}),
).AddHandlers(
	new(AddReactComponent),
).BuildCommand()

type AddReactComponent struct {
	tasks []task.Task
}

func (arc *AddReactComponent) Handle(c *cli.Context, cfg interface{}) error {
	rcc := cfg.(*ReactComponentConfig)
	fileName := strcase.ToKebab(rcc.Name)
	varName := strcase.ToCamel(rcc.Name)
	outDir := rcc.OutDir

	if rcc.AddFolder {
		outDir = path.Join(outDir, fileName)
		fileName = "index"
	}
	arc.addComponent(varName, fileName, outDir)

	if rcc.AddStory {
		arc.addStory(varName, fileName, outDir)
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

func (arc *AddReactComponent) addComponent(varName string, fileName string, outDir string) {
	b := new(bytes.Buffer)
	arc.tasks = append(arc.tasks,
		task.NewSequentialScheduler().
			AddTask(task.NewReadFromProviderTask(provider.NewEmbedFileProvider("templates/react/react_component.gotmpl"), b)).
			AddTask(task.NewTemplateRenderTask("react_component", map[string]string{"name": varName}, b, b)).
			AddTask(task.NewWriteFileToDiskTask(fileName, file.ReactComponentSuffix, outDir, false, b)))
}

func (arc *AddReactComponent) addStory(componentName string, fileName string, outDir string) {
	b := new(bytes.Buffer)
	arc.tasks = append(arc.tasks,
		task.NewSequentialScheduler().
			AddTask(task.NewReadFromProviderTask(provider.NewEmbedFileProvider("templates/react/react_storybook.gotmpl"), b)).
			AddTask(task.NewTemplateRenderTask(
				"react_storybook",
				map[string]string{"name": componentName, "file": fileName},
				b, b)).
			AddTask(task.NewWriteFileToDiskTask(fileName, file.StorybookSuffix, outDir, false, b)))
}

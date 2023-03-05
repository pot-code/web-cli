package cmd

import (
	"bytes"
	"path"

	"github.com/iancoleman/strcase"
	"github.com/pot-code/web-cli/pkg/command"
	"github.com/pot-code/web-cli/pkg/provider"
	"github.com/pot-code/web-cli/pkg/task"
	"github.com/urfave/cli/v2"
)

const (
	StorybookSuffix      = ".stories.tsx"
	ReactTestSuffix      = ".test.tsx"
	ReactComponentSuffix = ".tsx"
)

type ReactComponentConfig struct {
	AddTest   bool   `flag:"add-test" alias:"t" usage:"add test file"`
	AddStory  bool   `flag:"add-storybook" alias:"s" usage:"add storybook"`
	AddFolder bool   `flag:"add-folder" alias:"f" usage:"generate files in a folder with the name as the component"`
	OutDir    string `flag:"output" alias:"o" usage:"destination directory"`
	Name      string `arg:"0" alias:"COMPONENT_NAME" validate:"required,nature"`
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
	config := cfg.(*ReactComponentConfig)
	cn := strcase.ToCamel(config.Name)

	outDir := config.OutDir
	if config.AddFolder {
		outDir = path.Join(outDir, cn)
	}

	arc.addComponent(cn, outDir)
	if config.AddStory {
		arc.addStory(cn, outDir)
	}
	if config.AddTest {
		arc.addTest(cn, outDir)
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

func (arc *AddReactComponent) addComponent(componentName string, outDir string) {
	b := new(bytes.Buffer)
	arc.tasks = append(arc.tasks,
		task.NewSequentialScheduler().
			AddTask(task.NewReadFromProviderTask(provider.NewEmbedFileProvider("templates/react_component.tmpl"), b)).
			AddTask(task.NewTemplateRenderTask("react_component", map[string]string{"name": componentName}, b, b)).
			AddTask(task.NewWriteFileToDiskTask(componentName, ReactComponentSuffix, outDir, false, b)))
}

func (arc *AddReactComponent) addStory(componentName string, outDir string) {
	b := new(bytes.Buffer)
	arc.tasks = append(arc.tasks,
		task.NewSequentialScheduler().
			AddTask(task.NewReadFromProviderTask(provider.NewEmbedFileProvider("templates/react_storybook.tmpl"), b)).
			AddTask(task.NewTemplateRenderTask("react_storybook", map[string]string{"name": componentName}, b, b)).
			AddTask(task.NewWriteFileToDiskTask(componentName, StorybookSuffix, outDir, false, b)))
}

func (arc *AddReactComponent) addTest(componentName string, outDir string) {
	b := new(bytes.Buffer)
	arc.tasks = append(arc.tasks,
		task.NewSequentialScheduler().
			AddTask(task.NewReadFromProviderTask(provider.NewEmbedFileProvider("templates/react_test.tmpl"), b)).
			AddTask(task.NewTemplateRenderTask("react_test", map[string]string{"name": componentName}, b, b)).
			AddTask(task.NewWriteFileToDiskTask(componentName, ReactTestSuffix, outDir, false, b)))
}

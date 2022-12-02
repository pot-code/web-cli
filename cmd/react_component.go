package cmd

import (
	"fmt"
	"path"

	"github.com/iancoleman/strcase"
	"github.com/pot-code/web-cli/internal/command"
	"github.com/pot-code/web-cli/internal/task"
	"github.com/pot-code/web-cli/internal/template"
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
		return fmt.Errorf("execute ParallelExecutor: %w", err)
	}
	return nil
}

func (arc *AddReactComponent) addComponent(componentName string, outDir string) {
	arc.tasks = append(arc.tasks, task.NewGenerateFileFromTemplateTask(
		componentName,
		ReactComponentSuffix,
		outDir,
		false,
		componentName,
		template.NewLocalTemplateProvider(GetAbsoluteTemplatePath("react_component.tmpl")),
		map[string]string{
			"name": componentName,
		},
	))
}

func (arc *AddReactComponent) addStory(componentName string, outDir string) {
	arc.tasks = append(arc.tasks, task.NewGenerateFileFromTemplateTask(
		componentName,
		StorybookSuffix,
		outDir,
		false,
		componentName,
		template.NewLocalTemplateProvider(GetAbsoluteTemplatePath("react_storybook.tmpl")),
		map[string]string{
			"name": componentName,
		},
	))
}

func (arc *AddReactComponent) addTest(componentName string, outDir string) {
	arc.tasks = append(arc.tasks, task.NewGenerateFileFromTemplateTask(
		componentName,
		ReactTestSuffix,
		outDir,
		false,
		componentName,
		template.NewLocalTemplateProvider(GetAbsoluteTemplatePath("react_test.tmpl")),
		map[string]string{
			"name": componentName,
		},
	))
}

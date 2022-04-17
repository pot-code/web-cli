package cmd

import (
	"bytes"
	"path"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/pot-code/web-cli/internal/command"
	"github.com/pot-code/web-cli/internal/task"
	"github.com/pot-code/web-cli/templates"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type ReactComponentConfig struct {
	Test   bool   `flag:"test" alias:"t" usage:"add test file"`
	Story  bool   `flag:"storybook" alias:"sb" usage:"add storybook"`
	OutDir string `flag:"output" alias:"o" usage:"output directory"`
	Name   string `arg:"0" alias:"component_name" validate:"required,var"`
}

var ReactComponentCmd = command.NewCliCommand("component", "add react component",
	new(ReactComponentConfig),
	command.WithArgUsage("component_name"),
	command.WithAlias([]string{"c"}),
).AddHandlers(
	new(AddReactComponent),
).BuildCommand()

type AddReactComponent struct {
	componentName string
}

func (arc *AddReactComponent) Handle(c *cli.Context, cfg interface{}) error {
	config := cfg.(*ReactComponentConfig)
	componentName := strcase.ToCamel(config.Name)
	outDir := arc.getOutputPath(config.OutDir, componentName)
	tree := task.NewFileGenerationTree(outDir)

	arc.componentName = componentName
	log.WithFields(log.Fields{"handler": "AddReactComponent", "path": outDir}).Debug("output path")
	tree.AddNodes(arc.component())

	if config.Story {
		tree.AddNodes(arc.story())
	}

	if config.Test {
		tree.AddNodes(arc.test())
	}

	return task.NewParallelExecutor(
		task.BatchFileGenerationTask(
			tree.Flatten(),
		),
	).Run()
}

func (arc *AddReactComponent) component() *task.FileGenerator {
	name := arc.componentName
	return &task.FileGenerator{
		Name: getComponentFileName(name),
		Data: bytes.NewBufferString(templates.ReactComponent(name)),
	}
}

func (arc *AddReactComponent) story() *task.FileGenerator {
	name := arc.componentName
	return &task.FileGenerator{
		Name: getStoryFileName(name),
		Data: bytes.NewBufferString(templates.ReactStory(name)),
	}
}

func (arc *AddReactComponent) test() *task.FileGenerator {
	name := arc.componentName
	return &task.FileGenerator{
		Name: getTestFileName(name),
		Data: bytes.NewBufferString(templates.ReactTest(name)),
	}
}

func (arc *AddReactComponent) getOutputPath(dir, componentName string) string {
	if strings.HasSuffix(dir, componentName) {
		return dir
	}
	return path.Join(dir, componentName)
}

func getStoryFileName(componentName string) string {
	return componentName + ".stories.tsx"
}

func getTestFileName(componentName string) string {
	return componentName + ".test.tsx"
}

func getComponentFileName(componentName string) string {
	return componentName + ".tsx"
}

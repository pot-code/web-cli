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
	Story  bool   `flag:"storybook" alias:"s" usage:"add storybook"`
	Folder bool   `flag:"folder" alias:"f" usage:"generate files in a folder with the name as the component"`
	OutDir string `flag:"output" alias:"o" usage:"destination directory"`
	Name   string `arg:"0" alias:"COMPONENT_NAME" validate:"required,nature"`
}

var ReactComponentCmd = command.NewCliCommand("component", "add react component",
	new(ReactComponentConfig),
	command.WithArgUsage("COMPONENT_NAME"),
	command.WithAlias([]string{"c"}),
).AddHandlers(
	new(AddReactComponent),
).BuildCommand()

type AddReactComponent struct{}

func (arc *AddReactComponent) Handle(c *cli.Context, cfg interface{}) error {
	config := cfg.(*ReactComponentConfig)
	componentName := strcase.ToCamel(config.Name)
	outDir := config.OutDir
	if config.Folder {
		outDir = getIsolatedOutputPath(config.OutDir, componentName)
	}

	tree := task.NewFileGenerationTree(outDir)
	log.WithFields(log.Fields{"handler": "AddReactComponent", "path": outDir}).Debug("output path")
	tree.AddNodes(arc.component(componentName))
	if config.Story {
		tree.AddNodes(arc.story(componentName))
	}
	if config.Test {
		tree.AddNodes(arc.test(componentName))
	}
	return task.NewParallelExecutor(
		task.BatchFileGenerationTask(
			tree.Flatten(),
		),
	).Run()
}

func (arc *AddReactComponent) component(componentName string) *task.FileGenerator {
	return &task.FileGenerator{
		Name: getComponentFileName(componentName),
		Data: bytes.NewBufferString(templates.ReactComponent(componentName)),
	}
}

func (arc *AddReactComponent) story(componentName string) *task.FileGenerator {
	return &task.FileGenerator{
		Name: getStoryFileName(componentName),
		Data: bytes.NewBufferString(templates.ReactStory(componentName)),
	}
}

func (arc *AddReactComponent) test(componentName string) *task.FileGenerator {
	return &task.FileGenerator{
		Name: getTestFileName(componentName),
		Data: bytes.NewBufferString(templates.ReactTest(componentName)),
	}
}

func getIsolatedOutputPath(dir, componentName string) string {
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

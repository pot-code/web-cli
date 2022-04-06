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
	Scss   bool   `flag:"scss" alias:"s" usage:"add scss module"`
	Story  bool   `flag:"storybook" alias:"sb" usage:"add storybook"`
	OutDir string `flag:"dir" alias:"o" usage:"output directory"`
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

	var scss string
	if config.Scss {
		s := arc.scss()
		scss = s.Name
		tree.AddNodes(s)
	}
	tree.AddNodes(arc.component(scss))

	if config.Story {
		tree.AddNodes(arc.story())
	}

	return task.NewParallelExecutor(
		task.BatchFileGenerationTask(
			tree.Flatten(),
		),
	).Run()
}

func (arc *AddReactComponent) component(scss string) *task.FileGenerator {
	name := arc.componentName
	return &task.FileGenerator{
		Name: arc.getComponentFileName(),
		Data: bytes.NewBufferString(templates.ReactComponent(name, scss)),
	}
}

func (arc *AddReactComponent) scss() *task.FileGenerator {
	rootClass := strcase.ToKebab(arc.componentName)
	return &task.FileGenerator{
		Name: arc.getScssFileName(),
		Data: bytes.NewBufferString(templates.ReactSCSS(rootClass)),
	}
}

func (arc *AddReactComponent) story() *task.FileGenerator {
	name := arc.componentName
	return &task.FileGenerator{
		Name: arc.getStoryFileName(),
		Data: bytes.NewBufferString(templates.ReactStory(name)),
	}
}

func (arc *AddReactComponent) getOutputPath(dir, componentName string) string {
	if strings.HasSuffix(dir, componentName) {
		return dir
	}
	return path.Join(dir, componentName)
}

func (arc *AddReactComponent) getStoryFileName() string {
	return arc.componentName + ".stories.tsx"
}

func (arc *AddReactComponent) getScssFileName() string {
	return arc.componentName + ".module.scss"
}

func (arc *AddReactComponent) getComponentFileName() string {
	return arc.componentName + ".tsx"
}

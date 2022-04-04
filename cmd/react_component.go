package cmd

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/pot-code/web-cli/internal/command"
	"github.com/pot-code/web-cli/internal/constant"
	"github.com/pot-code/web-cli/internal/shell"
	"github.com/pot-code/web-cli/internal/task"
	"github.com/pot-code/web-cli/templates"
	"github.com/urfave/cli/v2"
)

type ReactComponentConfig struct {
	Hook    bool   `flag:"hook" usage:"add react hook"`
	Scss    bool   `flag:"scss" alias:"S" usage:"add scss module"`
	Story   bool   `flag:"storybook" alias:"sb" usage:"add storybook"`
	Emotion bool   `flag:"emotion" alias:"e" usage:"add @emotion/react"`
	Dir     string `flag:"dir" alias:"d" usage:"output dir"`
	Name    string `arg:"0" alias:"component_name" validate:"required,var"`
}

var ReactComponentCmd = command.NewCliCommand("component", "add react component",
	new(ReactComponentConfig),
	command.WithArgUsage("component_name"),
	command.WithAlias([]string{"c"}),
).AddFeature(
	new(AddReactComponent),
	new(AddReactEmotionFeat),
).ExportCommand()

type AddReactComponent struct {
	ComponentName string
}

var _ command.CommandFeature = &AddReactComponent{}

func (arc *AddReactComponent) Cond(c *cli.Context) bool {
	return true
}

func (arc *AddReactComponent) Handle(c *cli.Context, cfg interface{}) error {
	config := cfg.(*ReactComponentConfig)
	name := config.Name

	arc.ComponentName = strings.ReplaceAll(name, "-", "_")
	return arc.addReactComponent(config).Run()
}

func (arc *AddReactComponent) addReactComponent(cfg *ReactComponentConfig) task.Task {
	// dir := cfg.Dir
	name := arc.ComponentName

	var styleFile string
	var files = []*task.FileGenerator{
		{
			Name: arc.getComponentFileName(),
			Data: bytes.NewBufferString(templates.ReactComponent(name, styleFile)),
		},
	}

	if cfg.Scss {
		styleFile = arc.getScssFileName()
		files = append(files, arc.addScss(cfg))
	}

	if cfg.Story {
		files = append(files, arc.addStoryBook(cfg))
	}

	return task.NewParallelExecutor(task.BatchFileGenerationTask(files))
}

func (arc *AddReactComponent) addScss(cfg *ReactComponentConfig) *task.FileGenerator {
	rootClass := strcase.ToKebab(arc.ComponentName)

	return &task.FileGenerator{
		Name: arc.getScssFileName(),
		Data: bytes.NewBufferString(templates.ReactSCSS(rootClass)),
	}
}

func (arc *AddReactComponent) getScssFileName() string {
	return fmt.Sprintf(constant.ReactScssPattern, arc.ComponentName)
}

func (arc *AddReactComponent) addStoryBook(cfg *ReactComponentConfig) *task.FileGenerator {
	name := arc.ComponentName

	return &task.FileGenerator{
		Name: arc.getStoryFileName(),
		Data: bytes.NewBufferString(templates.ReactStory(name)),
	}
}

func (arc *AddReactComponent) getStoryFileName() string {
	return fmt.Sprintf(constant.ReactStorybookPattern, arc.ComponentName, constant.TsxSuffix)
}

func (arc *AddReactComponent) getComponentFileName() string {
	return fmt.Sprintf(constant.ReactComponentPattern, arc.ComponentName, constant.TsxSuffix)
}

type AddReactEmotionFeat struct{}

var _ command.CommandFeature = &AddReactEmotionFeat{}

func (arc *AddReactEmotionFeat) Cond(c *cli.Context) bool {
	return c.Bool("emotion")
}

func (arc *AddReactEmotionFeat) Handle(c *cli.Context, cfg interface{}) error {
	return task.NewParallelExecutor(
		[]task.Task{
			&task.FileGenerator{
				Name: ".babelrc",
				Data: bytes.NewBufferString(templates.ReactEmotion()),
			},
			shell.YarnAdd("@emotion/react"),
			shell.YarnAddDev("@emotion/babel-preset-css-prop"),
		},
	).Run()
}

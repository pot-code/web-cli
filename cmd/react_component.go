package cmd

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/pot-code/web-cli/internal/commands"
	"github.com/pot-code/web-cli/internal/constants"
	"github.com/pot-code/web-cli/internal/task"
	"github.com/pot-code/web-cli/internal/util"
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

var ReactComponentCmd = util.NewCliCommand("component", "add react component",
	new(ReactComponentConfig),
	util.WithArgUsage("component_name"),
	util.WithAlias([]string{"c"}),
).AddFeature(
	new(AddReactComponent),
	new(AddReactEmotionFeat),
).ExportCommand()

type AddReactComponent struct {
	ComponentName string
}

var _ util.CommandFeature = &AddReactComponent{}

func (arc *AddReactComponent) Cond(c *cli.Context) bool {
	return true
}

func (arc *AddReactComponent) Handle(c *cli.Context, cfg interface{}) error {
	config := cfg.(*ReactComponentConfig)
	name := config.Name

	arc.ComponentName = strings.ReplaceAll(name, "-", "_")
	return arc.addReactComponent(config).Run()
}

func (arc *AddReactComponent) addReactComponent(cfg *ReactComponentConfig) task.Runner {
	dir := cfg.Dir
	name := arc.ComponentName

	var styleFile string
	var files = []*task.FileDesc{
		{
			Path: arc.getComponentFileName(),
			Source: func(buf *bytes.Buffer) error {
				templates.WriteReactComponent(buf, name, styleFile)
				return nil
			},
		},
	}

	if cfg.Scss {
		styleFile = arc.getScssFileName()
		files = append(files, arc.addScss(cfg))
	}

	if cfg.Story {
		files = append(files, arc.addStoryBook(cfg))
	}

	return util.NewTaskComposer(dir).AddFile(files...)
}

func (arc *AddReactComponent) addScss(cfg *ReactComponentConfig) *task.FileDesc {
	rootClass := strcase.ToKebab(arc.ComponentName)

	return &task.FileDesc{
		Path: arc.getScssFileName(),
		Source: func(buf *bytes.Buffer) error {
			templates.WriteReactSCSS(buf, rootClass)
			return nil
		},
	}
}

func (arc *AddReactComponent) getScssFileName() string {
	return fmt.Sprintf(constants.ReactScssPattern, arc.ComponentName)
}

func (arc *AddReactComponent) addStoryBook(cfg *ReactComponentConfig) *task.FileDesc {
	name := arc.ComponentName

	return &task.FileDesc{
		Path: arc.getStoryFileName(),
		Source: func(buf *bytes.Buffer) error {
			templates.WriteReactStory(buf, name)
			return nil
		},
	}
}

func (arc *AddReactComponent) getStoryFileName() string {
	return fmt.Sprintf(constants.ReactStorybookPattern, arc.ComponentName, constants.TsxSuffix)
}

func (arc *AddReactComponent) getComponentFileName() string {
	return fmt.Sprintf(constants.ReactComponentPattern, arc.ComponentName, constants.TsxSuffix)
}

type AddReactEmotionFeat struct{}

var _ util.CommandFeature = &AddReactEmotionFeat{}

func (arc *AddReactEmotionFeat) Cond(c *cli.Context) bool {
	return c.Bool("emotion")
}

func (arc *AddReactEmotionFeat) Handle(c *cli.Context, cfg interface{}) error {
	return util.NewTaskComposer("").AddFile(
		&task.FileDesc{
			Path: ".babelrc",
			Source: func(buf *bytes.Buffer) error {
				templates.WriteReactEmotion(buf)
				return nil
			},
		},
	).AddCommand(
		commands.YarnAdd("@emotion/react"),
		commands.YarnAddDev("@emotion/babel-preset-css-prop"),
	).Run()
}

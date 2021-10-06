package cmd

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/pot-code/web-cli/pkg/commands"
	"github.com/pot-code/web-cli/pkg/constants"
	"github.com/pot-code/web-cli/pkg/core"
	"github.com/pot-code/web-cli/pkg/util"
	"github.com/pot-code/web-cli/templates"
	"github.com/urfave/cli/v2"
)

type ReactComponentConfig struct {
	Hook    bool   `flag:"hook" usage:"add react hook"`
	Scss    bool   `flag:"scss" alias:"S" usage:"add scss module"`
	Story   bool   `flag:"story" alias:"sb" usage:"add storybook"`
	Emotion bool   `flag:"emotion" alias:"e" usage:"add @emotion/react"`
	Dir     string `flag:"dir" alias:"d" usage:"output dir"`
	Name    string `arg:"0" alias:"component_name" validate:"required,var"`
}

var ReactComponentCmd = core.NewCliLeafCommand("component", "add react component",
	new(ReactComponentConfig),
	core.WithArgUsage("component_name"),
	core.WithAlias([]string{"c"}),
).AddService(
	new(AddReactComponentService),
	new(AddReactEmotionService),
).ExportCommand()

type AddReactComponentService struct {
	ComponentName string
}

var _ core.CommandService = &AddReactComponentService{}

func (arc *AddReactComponentService) Cond(c *cli.Context) bool {
	return true
}

func (arc *AddReactComponentService) Handle(c *cli.Context, cfg interface{}) error {
	config := cfg.(*ReactComponentConfig)
	name := config.Name

	arc.ComponentName = strings.ReplaceAll(name, "-", "_")
	return arc.addReactComponent(config).Run()
}

func (arc *AddReactComponentService) addScss(cfg *ReactComponentConfig) *core.FileDesc {
	rootClass := strcase.ToKebab(arc.ComponentName)

	return &core.FileDesc{
		Path: arc.getScssFileName(),
		Source: func(buf *bytes.Buffer) error {
			templates.WriteReactSCSS(buf, rootClass)
			return nil
		},
	}
}

func (arc *AddReactComponentService) addStoryBook(cfg *ReactComponentConfig) *core.FileDesc {
	name := arc.ComponentName

	return &core.FileDesc{
		Path: arc.getStoryFileName(),
		Source: func(buf *bytes.Buffer) error {
			templates.WriteReactStory(buf, name)
			return nil
		},
	}
}

func (arc *AddReactComponentService) addReactComponent(cfg *ReactComponentConfig) core.Runner {
	dir := cfg.Dir
	name := arc.ComponentName

	var styleFile string
	var files = []*core.FileDesc{
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

func (arc *AddReactComponentService) getScssFileName() string {
	return fmt.Sprintf(constants.ReactScssPattern, arc.ComponentName)
}

func (arc *AddReactComponentService) getStoryFileName() string {
	return fmt.Sprintf(constants.ReactStorybookPattern, arc.ComponentName, constants.TsxSuffix)
}

func (arc *AddReactComponentService) getComponentFileName() string {
	return fmt.Sprintf(constants.ReactComponentPattern, arc.ComponentName, constants.TsxSuffix)
}

type AddReactEmotionService struct{}

var _ core.CommandService = &AddReactEmotionService{}

func (arc *AddReactEmotionService) Cond(c *cli.Context) bool {
	return c.Bool("emotion")
}

func (arc *AddReactEmotionService) Handle(c *cli.Context, cfg interface{}) error {
	return util.NewTaskComposer("").AddFile(
		&core.FileDesc{
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

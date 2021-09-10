package cmd

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/pot-code/web-cli/pkg/core"
	"github.com/pot-code/web-cli/pkg/util"
	"github.com/pot-code/web-cli/templates"
	"github.com/urfave/cli/v2"
)

type AddReactConfig struct {
	Hook    bool   `flag:"hook" usage:"add react hook"`
	Scss    bool   `flag:"scss" alias:"S" usage:"add scss module"`
	Story   bool   `flag:"story" alias:"sb" usage:"add storybook"`
	Emotion bool   `flag:"emotion" alias:"e" usage:"add @emotion/react"`
	Dir     string `flag:"dir" alias:"d" usage:"output dir"`
	Name    string `arg:"0" alias:"component_name" validate:"required,var"`
}

var AddReactCmd = core.NewCliLeafCommand("react", "add React components", new(AddReactConfig),
	core.WithArgUsage("component_name"),
	core.WithAlias([]string{"r"}),
).AddService(new(AddReactComponentService)).AddService(new(AddReactEmotionService)).ExportCommand()

type AddReactComponentService struct {
	ComponentName string
}

var _ core.CommandService = &AddReactComponentService{}

func (arc *AddReactComponentService) Cond(c *cli.Context) bool {
	return true
}

func (arc *AddReactComponentService) Handle(c *cli.Context, cfg interface{}) error {
	config := cfg.(*AddReactConfig)
	name := config.Name

	arc.ComponentName = strings.ReplaceAll(name, "-", "_")
	cmd := arc.addReactComponent(config)
	return cmd.Run()
}

func (arc *AddReactComponentService) addScss(cfg *AddReactConfig) *core.FileDesc {
	rootClass := strcase.ToKebab(arc.ComponentName)

	return &core.FileDesc{
		Path: arc.getScssFileName(),
		Data: func() []byte {
			var buf bytes.Buffer

			templates.WriteReactSCSS(&buf, rootClass)
			return buf.Bytes()
		},
	}
}

func (arc *AddReactComponentService) addStoryBook(cfg *AddReactConfig) *core.FileDesc {
	name := arc.ComponentName

	return &core.FileDesc{
		Path: arc.getStoryFileName(),
		Data: func() []byte {
			var buf bytes.Buffer

			templates.WriteReactStory(&buf, name)
			return buf.Bytes()
		},
	}
}

func (arc *AddReactComponentService) addReactComponent(cfg *AddReactConfig) core.Runner {
	dir := cfg.Dir
	name := arc.ComponentName

	var styleFile string
	var files = []*core.FileDesc{
		{
			Path: arc.getComponentFileName(),
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteReactComponent(&buf, name, styleFile)
				return buf.Bytes()
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
	return fmt.Sprintf("%s.%s.%s", arc.ComponentName, "module", "scss")
}

func (arc *AddReactComponentService) getStoryFileName() string {
	return fmt.Sprintf("%s.%s.%s", arc.ComponentName, "stories", "tsx")
}

func (arc *AddReactComponentService) getComponentFileName() string {
	return fmt.Sprintf("%s.%s", arc.ComponentName, "tsx")
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
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteReactEmotion(&buf)
				return buf.Bytes()
			},
		},
	).AddCommand(
		&core.Command{
			Bin:  "npm",
			Args: []string{"i", "@emotion/react"},
		},
		&core.Command{
			Bin:  "npm",
			Args: []string{"i", "-D", "@emotion/babel-preset-css-prop"},
		},
	).Run()
}

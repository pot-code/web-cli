package cmd

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/pot-code/web-cli/pkg/core"
	"github.com/pot-code/web-cli/pkg/util"
	"github.com/pot-code/web-cli/templates"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type addReactConfig struct {
	Hook  bool   `flag:"hook"`
	Scss  bool   `flag:"scss"`
	Story bool   `flag:"story"`
	Dir   string `flag:"dir"`
	Name  string `arg:"0" alias:"component_name" validate:"required,var"`
}

var addReactCmd = &cli.Command{
	Name:      "react",
	Usage:     "add React components",
	Aliases:   []string{"r"},
	ArgsUsage: "component_name",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "scss",
			Aliases: []string{"s"},
			Usage:   "add scss module",
			Value:   false,
		},
		&cli.BoolFlag{
			Name:    "story",
			Aliases: []string{"sb"},
			Usage:   "add scss module",
			Value:   false,
		},
		&cli.StringFlag{
			Name:    "dir",
			Aliases: []string{"d"},
			Usage:   "output dir",
			Value:   "",
		},
		&cli.BoolFlag{
			Name:    "emotion",
			Aliases: []string{"e"},
			Usage:   "add @emotion/react",
			Value:   false,
		},
	},
	Action: func(c *cli.Context) error {
		if c.Bool("emotion") {
			cmd := addReactEmotion()
			return cmd.Run()
		}

		config := new(addReactConfig)
		err := util.ParseConfig(c, config)
		if err != nil {
			if _, ok := err.(*util.CommandError); ok {
				cli.ShowCommandHelp(c, c.Command.Name)
			}
			return err
		}

		name := config.Name
		name = strings.ReplaceAll(name, "-", "_")
		log.WithFields(log.Fields{
			"cmd": "addReactCmd",
		}).Debug("proceed component name: ", name)

		cmd := addReactComponent(strcase.ToCamel(name), config.Dir, config.Scss, config.Story)

		return cmd.Run()
	},
}

func addReactEmotion() core.Runner {
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
	)
}

func addReactComponent(name, dir string, scss, story bool) core.Runner {
	var (
		stylePath string
		desc      []*core.FileDesc
	)

	if scss {
		styleName := strcase.ToKebab(name)
		stylePath = fmt.Sprintf("%s.%s.%s", name, "module", "scss")
		desc = append(desc, &core.FileDesc{
			Path: stylePath,
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteReactSCSS(&buf, styleName)
				return buf.Bytes()
			},
		})
	}

	if story {
		desc = append(desc, &core.FileDesc{
			Path: fmt.Sprintf("%s.%s.%s", name, "stories", "tsx"),
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteReactStory(&buf, name)
				return buf.Bytes()
			},
		})
	}

	desc = append(desc,
		&core.FileDesc{
			Path: fmt.Sprintf("%s.%s", name, "tsx"),
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteReactComponent(&buf, name, stylePath)
				return buf.Bytes()
			},
		})

	return util.NewTaskComposer(dir).AddFile(desc...)
}

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

const cmdReactName = "react"

type addReactConfig struct {
	Hook  bool   `name:"hook"`
	Style bool   `name:"style"`
	Name  string `arg:"0" name:"NAME" validate:"required,var"`
}

var addReactCmd = &cli.Command{
	Name:      cmdReactName,
	Usage:     "add React components",
	Aliases:   []string{"r"},
	ArgsUsage: "NAME",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "hook",
			Aliases: []string{"H"},
			Usage:   "add hook",
			Value:   false,
		},
		&cli.BoolFlag{
			Name:    "style",
			Aliases: []string{"s"},
			Usage:   "add scss module style",
			Value:   false,
		},
	},
	Action: func(c *cli.Context) error {
		config := new(addReactConfig)
		err := util.ParseConfig(c, config)
		if err != nil {
			if _, ok := err.(*util.CommandError); ok {
				cli.ShowCommandHelp(c, cmdReactName)
			}
			return err
		}

		name := config.Name
		name = strings.ReplaceAll(name, "-", "_")
		log.Debug("preprocessed component name: ", name)

		var cmd core.Generator
		if config.Hook {
			if !strings.HasPrefix(name, "use") {
				name = "use" + strcase.ToCamel(name)
			}
			cmd = addReactHook(strcase.ToLowerCamel(name))
		} else {
			cmd = addReactComponent(strcase.ToCamel(name), config.Style)
		}

		err = cmd.Run()
		if err != nil {
			cmd.Cleanup()
		}
		return err
	},
}

func addReactHook(name string) core.Generator {
	return util.NewTaskComposer("",
		&core.FileDesc{
			Path: fmt.Sprintf("%s.%s", name, "ts"),
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteReactHook(&buf, name)
				return buf.Bytes()
			},
		},
	)
}

func addReactComponent(name string, style bool) core.Generator {
	var (
		stylePath string
		desc      []*core.FileDesc
	)
	if style {
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
	desc = append(desc,
		&core.FileDesc{
			Path: fmt.Sprintf("%s.%s", name, "tsx"),
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteReactComponent(&buf, name, stylePath)
				return buf.Bytes()
			},
		})
	return util.NewTaskComposer("", desc...)
}

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
	Hook bool   `name:"hook"`
	Name string `arg:"0" name:"NAME" validate:"required,var"`
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
			cmd = newAddReactHook(strcase.ToLowerCamel(name))
		} else {
			cmd = newAddReactComponent(strcase.ToCamel(name))
		}

		err = cmd.Run()
		if err != nil {
			cmd.Cleanup()
		}
		return err
	},
}

func newAddReactHook(name string) core.Generator {
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

func newAddReactComponent(name string) core.Generator {
	return util.NewTaskComposer("",
		&core.FileDesc{
			Path: fmt.Sprintf("%s.%s", name, "tsx"),
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteReactComponent(&buf, name)
				return buf.Bytes()
			},
		},
	)
}

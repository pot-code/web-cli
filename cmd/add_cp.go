package cmd

import (
	"bytes"

	"github.com/pot-code/web-cli/pkg/core"
	"github.com/pot-code/web-cli/pkg/util"
	"github.com/pot-code/web-cli/templates"
	"github.com/urfave/cli/v2"
)

type addPriceUpdateConfig struct {
	Name string `arg:"0" name:"NAME" validate:"required"`
}

var addPriceUpdateConfigCmd = &cli.Command{
	Name:      "cp",
	Usage:     "add a price update config folder",
	ArgsUsage: "NAME",
	Action: func(c *cli.Context) error {
		config := new(addPriceUpdateConfig)
		err := util.ParseConfig(c, config)
		if err != nil {
			if _, ok := err.(*util.CommandError); ok {
				cli.ShowCommandHelp(c, c.Command.Name)
			}
			return err
		}

		cmd := addPriceUpdateEntry(config.Name)
		return cmd.Run()
	},
}

func addPriceUpdateEntry(name string) core.Runner {
	return util.NewTaskComposer(name).AddFile(
		&core.FileDesc{
			Path: "provider.yml",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoProviders(&buf)
				return buf.Bytes()
			},
		},
		&core.FileDesc{
			Path: "consumers.yml",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoConsumers(&buf)
				return buf.Bytes()
			},
		},
	)
}

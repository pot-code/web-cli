package cmd

import (
	"bytes"

	"github.com/pot-code/web-cli/pkg/task"
	"github.com/pot-code/web-cli/pkg/util"
	"github.com/pot-code/web-cli/templates"
	"github.com/urfave/cli/v2"
)

type AddPriceUpdateConfig struct {
	Name string `arg:"0" alias:"config_name" validate:"required"`
}

var AddPriceUpdateConfigCmd = util.NewCliCommand("cp", "add a price update config folder",
	new(AddPriceUpdateConfig),
	util.WithArgUsage("config_name"),
).AddFeature(AddPriceUpdateConfigFile).ExportCommand()

var AddPriceUpdateConfigFile = util.NoCondFeature(func(c *cli.Context, cfg interface{}) error {
	config := cfg.(*AddPriceUpdateConfig)

	return util.NewTaskComposer(config.Name).AddFile(
		&task.FileDesc{
			Path: "provider.yml",
			Source: func(buf *bytes.Buffer) error {
				templates.WriteGoProviders(buf)
				return nil
			},
		},
		&task.FileDesc{
			Path: "consumers.yml",
			Source: func(buf *bytes.Buffer) error {
				templates.WriteGoConsumers(buf)
				return nil
			},
		},
	).Run()
})

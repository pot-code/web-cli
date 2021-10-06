package cmd

import (
	"bytes"

	"github.com/pot-code/web-cli/pkg/core"
	"github.com/pot-code/web-cli/pkg/util"
	"github.com/pot-code/web-cli/templates"
	"github.com/urfave/cli/v2"
)

type AddPriceUpdateConfig struct {
	Name string `arg:"0" alias:"config_name" validate:"required"`
}

var AddPriceUpdateConfigCmd = core.NewCliLeafCommand("cp", "add a price update config folder",
	new(AddPriceUpdateConfig),
	core.WithArgUsage("config_name"),
).AddService(AddPriceUpdateConfigService).ExportCommand()

var AddPriceUpdateConfigService = util.NoCondFunctionService(func(c *cli.Context, cfg interface{}) error {
	config := cfg.(*AddPriceUpdateConfig)

	return util.NewTaskComposer(config.Name).AddFile(
		&core.FileDesc{
			Path: "provider.yml",
			Source: func(buf *bytes.Buffer) error {
				templates.WriteGoProviders(buf)
				return nil
			},
		},
		&core.FileDesc{
			Path: "consumers.yml",
			Source: func(buf *bytes.Buffer) error {
				templates.WriteGoConsumers(buf)
				return nil
			},
		},
	).Run()
})

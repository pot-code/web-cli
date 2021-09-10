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
	).Run()
})

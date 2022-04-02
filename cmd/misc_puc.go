package cmd

import (
	"bytes"

	"github.com/pot-code/web-cli/internal/command"
	"github.com/pot-code/web-cli/internal/task"
	"github.com/pot-code/web-cli/internal/util"
	"github.com/pot-code/web-cli/templates"
	"github.com/urfave/cli/v2"
)

type AddPriceUpdateConfig struct {
	Name string `arg:"0" alias:"config_name" validate:"required"`
}

var AddPriceUpdateConfigCmd = command.NewCliCommand("cp", "add a price update config folder",
	new(AddPriceUpdateConfig),
	command.WithArgUsage("config_name"),
).AddFeature(AddPriceUpdateConfigFile).ExportCommand()

var AddPriceUpdateConfigFile = util.NoCondFeature(func(c *cli.Context, cfg interface{}) error {
	config := cfg.(*AddPriceUpdateConfig)

	return task.NewParallelExecutor(
		task.BatchFileTask(
			task.NewFileRequestTree(config.Name).AddNode(
				&task.FileRequest{
					Name: "provider.yml",
					Data: bytes.NewBufferString(templates.GoProviders()),
				},
				&task.FileRequest{
					Name: "consumers.yml",
					Data: bytes.NewBufferString(templates.GoConsumers()),
				},
			).Flatten(),
		)...,
	).Run()
})

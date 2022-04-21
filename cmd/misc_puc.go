package cmd

import (
	"bytes"

	"github.com/pot-code/web-cli/internal/command"
	"github.com/pot-code/web-cli/internal/task"
	"github.com/pot-code/web-cli/templates"
	"github.com/urfave/cli/v2"
)

type AddPriceUpdateConfig struct {
	Name string `arg:"0" alias:"config_name" validate:"required"`
}

var AddPriceUpdateCmd = command.NewCliCommand("price", "add a price update config folder",
	new(AddPriceUpdateConfig),
	command.WithArgUsage("config_name"),
	command.WithAlias([]string{"p"}),
).AddHandlers(AddPriceUpdateConfigFile).BuildCommand()

var AddPriceUpdateConfigFile = command.InlineHandler(func(c *cli.Context, cfg interface{}) error {
	config := cfg.(*AddPriceUpdateConfig)

	return task.NewParallelExecutor(
		task.BatchFileGenerationTask(
			task.NewFileGenerationTree(config.Name).AddNodes(
				&task.FileGenerator{
					Name: "provider.yml",
					Data: bytes.NewBufferString(templates.GoProviders()),
				},
				&task.FileGenerator{
					Name: "consumers.yml",
					Data: bytes.NewBufferString(templates.GoConsumers()),
				},
			).Flatten(),
		),
	).Run()
})

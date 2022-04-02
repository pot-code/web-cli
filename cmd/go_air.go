package cmd

import (
	"bytes"

	"github.com/pot-code/web-cli/internal/command"
	"github.com/pot-code/web-cli/internal/task"
	"github.com/pot-code/web-cli/internal/util"
	"github.com/pot-code/web-cli/templates"
	"github.com/urfave/cli/v2"
)

var GoAirCmd = command.NewCliCommand("air", "add air live reload config", nil).
	AddFeature(AddAirConfig).ExportCommand()

var AddAirConfig = util.NoCondFeature(func(c *cli.Context, cfg interface{}) error {
	return task.NewParallelExecutor(
		&task.FileGenerator{
			Name: "air.toml",
			Data: bytes.NewBufferString(templates.GoAirConfig()),
		},
	).Run()
})

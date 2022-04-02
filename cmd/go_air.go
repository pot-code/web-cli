package cmd

import (
	"bytes"

	"github.com/pot-code/web-cli/internal/task"
	"github.com/pot-code/web-cli/internal/util"
	"github.com/pot-code/web-cli/templates"
	"github.com/urfave/cli/v2"
)

var GoAirCmd = util.NewCliCommand("air", "add air live reload config", nil).
	AddFeature(AddAirConfig).ExportCommand()

var AddAirConfig = util.NoCondFeature(func(c *cli.Context, cfg interface{}) error {
	return util.NewTaskComposer("").AddFile(
		&task.FileDesc{
			Path: "air.toml",
			Source: func(buf *bytes.Buffer) error {
				templates.WriteGoAirConfig(buf)
				return nil
			},
		},
	).Run()
})

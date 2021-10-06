package cmd

import (
	"bytes"

	"github.com/pot-code/web-cli/pkg/core"
	"github.com/pot-code/web-cli/pkg/util"
	"github.com/pot-code/web-cli/templates"
	"github.com/urfave/cli/v2"
)

var GoAirCmd = core.NewCliLeafCommand("air", "add air live reload config", nil).
	AddService(GoAirService).ExportCommand()

var GoAirService = util.NoCondFunctionService(func(c *cli.Context, cfg interface{}) error {
	return util.NewTaskComposer("").AddFile(
		&core.FileDesc{
			Path: "air.toml",
			Source: func(buf *bytes.Buffer) error {
				templates.WriteGoAirConfig(buf)
				return nil
			},
		},
	).Run()
})

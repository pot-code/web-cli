package cmd

import (
	"bytes"

	"github.com/pot-code/web-cli/pkg/core"
	"github.com/pot-code/web-cli/pkg/util"
	"github.com/pot-code/web-cli/templates"
	"github.com/urfave/cli/v2"
)

var AddGoAirCmd = core.NewCliLeafCommand("air", "add air live reload config", nil).
	AddService(AddGoAirService).ExportCommand()

var AddGoAirService = util.NoCondFunctionService(func(c *cli.Context, cfg interface{}) error {
	return util.NewTaskComposer("").AddFile(
		&core.FileDesc{
			Path: "air.toml",
			Data: func() ([]byte, error) {
				var buf bytes.Buffer

				templates.WriteGoAirConfig(&buf)
				return buf.Bytes(), nil
			},
		},
	).Run()
})

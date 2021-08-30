package cmd

import (
	"bytes"

	"github.com/pot-code/web-cli/pkg/core"
	"github.com/pot-code/web-cli/pkg/util"
	"github.com/pot-code/web-cli/templates"
	"github.com/urfave/cli/v2"
)

var addGoAirCmd = &cli.Command{
	Name:  "air",
	Usage: "add air live reload support",
	Action: func(c *cli.Context) error {
		cmd := addGoAir()
		err := cmd.Run()
		if err != nil {
			cmd.Cleanup()
		}
		return err
	},
}

func addGoAir() core.Generator {
	return util.NewTaskComposer("",
		&core.FileDesc{
			Path: "air.toml",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoAirConfig(&buf)
				return buf.Bytes()
			},
		},
	)
}

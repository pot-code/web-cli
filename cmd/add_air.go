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
		return cmd.Run()
	},
}

func addGoAir() core.Runner {
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

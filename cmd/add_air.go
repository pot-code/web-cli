package cmd

import (
	"bytes"

	"github.com/pot-code/web-cli/pkg/core"
	"github.com/pot-code/web-cli/pkg/util"
	"github.com/pot-code/web-cli/templates"
	"github.com/urfave/cli/v2"
)

var AddGoAirCmd = core.NewCliLeafCommand("air", "add air live reload support", nil).
	AddService(new(AddGoAirService)).ExportCommand()

type AddGoAirService struct{}

var _ core.CommandService = &AddGoAirService{}

func (aga *AddGoAirService) Cond(c *cli.Context) bool {
	return true
}

func (aga *AddGoAirService) Handle(c *cli.Context, cfg interface{}) error {
	return util.NewTaskComposer("").AddFile(
		&core.FileDesc{
			Path: "air.toml",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoAirConfig(&buf)
				return buf.Bytes()
			},
		},
	).Run()
}

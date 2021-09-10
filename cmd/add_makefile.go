package cmd

import (
	"bytes"

	"github.com/pot-code/web-cli/pkg/core"
	"github.com/pot-code/web-cli/pkg/util"
	"github.com/pot-code/web-cli/templates"
	"github.com/urfave/cli/v2"
)

var AddGoMakefileCmd = core.NewCliLeafCommand("makefile", "add Makefile", nil,
	core.WithAlias([]string{"m"}),
).AddService(new(AddGoMakefileService)).ExportCommand()

type AddGoMakefileService struct{}

var _ core.CommandService = &AddGoMakefileService{}

func (agm *AddGoMakefileService) Cond(c *cli.Context) bool {
	return true
}

func (agm *AddGoMakefileService) Handle(c *cli.Context, cfg interface{}) error {
	return util.NewTaskComposer("").AddFile(
		&core.FileDesc{
			Path: "Makefile",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoMakefile(&buf)
				return buf.Bytes()
			},
		},
	).Run()
}

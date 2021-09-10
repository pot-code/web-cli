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
).AddService(AddGoMakefileService).ExportCommand()

var AddGoMakefileService = util.NoCondFunctionService(func(c *cli.Context, cfg interface{}) error {
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
})

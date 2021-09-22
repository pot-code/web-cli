package cmd

import (
	"bytes"

	"github.com/pot-code/web-cli/pkg/core"
	"github.com/pot-code/web-cli/pkg/util"
	"github.com/pot-code/web-cli/templates"
	"github.com/urfave/cli/v2"
)

var GoMakefileCmd = core.NewCliLeafCommand("makefile", "add Makefile", nil,
	core.WithAlias([]string{"m"}),
).AddService(GoMakefileService).ExportCommand()

var GoMakefileService = util.NoCondFunctionService(func(c *cli.Context, cfg interface{}) error {
	return util.NewTaskComposer("").AddFile(
		&core.FileDesc{
			Path: "Makefile",
			Data: func() ([]byte, error) {
				var buf bytes.Buffer

				templates.WriteGoServerMakefile(&buf)
				return buf.Bytes(), nil
			},
		},
	).Run()
})

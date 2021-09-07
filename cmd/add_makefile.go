package cmd

import (
	"bytes"

	"github.com/pot-code/web-cli/pkg/core"
	"github.com/pot-code/web-cli/pkg/util"
	"github.com/pot-code/web-cli/templates"
	"github.com/urfave/cli/v2"
)

var addGoMakefileCmd = &cli.Command{
	Name:    "makefile",
	Aliases: []string{"m"},
	Usage:   "add Makefile",
	Action: func(c *cli.Context) error {
		cmd := addGoMakefile()
		return cmd.Run()
	},
}

func addGoMakefile() core.Runner {
	return util.NewTaskComposer("").AddFile(
		&core.FileDesc{
			Path: "Makefile",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoMakefile(&buf)
				return buf.Bytes()
			},
		},
	)
}

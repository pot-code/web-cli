package cmd

import (
	"bytes"

	"github.com/pot-code/web-cli/pkg/task"
	"github.com/pot-code/web-cli/pkg/util"
	"github.com/pot-code/web-cli/templates"
	"github.com/urfave/cli/v2"
)

var GoMakefileCmd = util.NewCliCommand("makefile", "add Makefile", nil,
	util.WithAlias([]string{"m"}),
).AddFeature(AddMakefile).ExportCommand()

var AddMakefile = util.NoCondFeature(func(c *cli.Context, cfg interface{}) error {
	return util.NewTaskComposer("").AddFile(
		&task.FileDesc{
			Path: "Makefile",
			Source: func(buf *bytes.Buffer) error {
				templates.WriteGoServerMakefile(buf)
				return nil
			},
		},
	).Run()
})

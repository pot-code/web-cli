package cmd

import (
	"bytes"

	"github.com/pot-code/web-cli/internal/command"
	"github.com/pot-code/web-cli/internal/task"
	"github.com/pot-code/web-cli/internal/util"
	"github.com/pot-code/web-cli/templates"
	"github.com/urfave/cli/v2"
)

var GoMakefileCmd = command.NewCliCommand("makefile", "add Makefile", nil,
	command.WithAlias([]string{"m"}),
).AddFeature(AddMakefile).ExportCommand()

var AddMakefile = util.NoCondFeature(func(c *cli.Context, cfg interface{}) error {
	return task.NewFileGenerator(
		&task.FileRequest{
			Name: "Makefile",
			Data:     bytes.NewBufferString(templates.GoServerMakefile()),
		},
	).Run()
})

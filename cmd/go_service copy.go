package cmd

import (
	"bytes"
	"path"

	"github.com/pkg/errors"
	"github.com/pot-code/web-cli/pkg/commands"
	"github.com/pot-code/web-cli/pkg/constants"
	"github.com/pot-code/web-cli/pkg/core"
	"github.com/pot-code/web-cli/pkg/transform"
	"github.com/pot-code/web-cli/pkg/util"
	"github.com/pot-code/web-cli/templates"
	"github.com/urfave/cli/v2"
)

type GoMigrateConfig struct {
}

var GoMigrateCmd = core.NewCliLeafCommand("migrate", "add migration",
	&GoMigrateConfig{},
	core.WithAlias([]string{"M"}),
).AddService(GoMigrateService).ExportCommand()

var GoMigrateService = util.NoCondFunctionService(func(c *cli.Context, cfg interface{}) error {
	meta, err := util.ParseGoMod(constants.GoModFile)
	if err != nil {
		return errors.WithStack(err)
	}

	return util.NewTaskComposer("").AddFile(
		&core.FileDesc{
			Path: path.Join("migrate", "config", "config.go"),
			Data: func() ([]byte, error) {
				var buf bytes.Buffer

				templates.WriteGoMigrateConfig(&buf)
				return buf.Bytes(), nil
			},
			Transforms: []core.Transform{transform.GoFormatSource},
		},
		&core.FileDesc{
			Path: path.Join("migrate", "migrate.go"),
			Data: func() ([]byte, error) {
				var buf bytes.Buffer

				templates.WriteGoMigrateMigration(&buf, meta.ProjectName, meta.Author)
				return buf.Bytes(), nil
			},
			Transforms: []core.Transform{transform.GoFormatSource},
		},
		&core.FileDesc{
			Path: path.Join("migrate", "wire.go"),
			Data: func() ([]byte, error) {
				var buf bytes.Buffer

				templates.WriteGoMigrateWire(&buf, meta.ProjectName, meta.Author)
				return buf.Bytes(), nil
			},
			Transforms: []core.Transform{transform.GoFormatSource},
		},
		&core.FileDesc{
			Path: path.Join("pkg", "db", "ent.go"),
			Data: func() ([]byte, error) {
				var buf bytes.Buffer

				templates.WriteGoMigratePkgEnt(&buf, meta.ProjectName, meta.Author)
				return buf.Bytes(), nil
			},
			Transforms: []core.Transform{transform.GoFormatSource},
		},
	).AddCommand(
		commands.GoModTidy(),
		commands.GoWire("./migrate"),
	).Run()
})

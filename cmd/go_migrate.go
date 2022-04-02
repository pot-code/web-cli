package cmd

import (
	"bytes"
	"path"

	"github.com/pkg/errors"
	"github.com/pot-code/web-cli/internal/commands"
	"github.com/pot-code/web-cli/internal/constants"
	"github.com/pot-code/web-cli/internal/task"
	"github.com/pot-code/web-cli/internal/transformation"
	"github.com/pot-code/web-cli/internal/util"
	"github.com/pot-code/web-cli/templates"
	"github.com/urfave/cli/v2"
)

type GoMigrateConfig struct{}

var GoMigrateCmd = util.NewCliCommand("migrate", "add migration",
	&GoMigrateConfig{},
	util.WithAlias([]string{"M"}),
).AddFeature(AddGoMigration).ExportCommand()

var AddGoMigration = util.NoCondFeature(func(c *cli.Context, cfg interface{}) error {
	meta, err := util.ParseGoMod(constants.GoModFile)
	if err != nil {
		return errors.WithStack(err)
	}

	found := false
	for _, r := range meta.Requires {
		if r.Mod.Path == "entgo.io/ent" {
			found = true
			break
		}
	}

	if !found {
		return errors.New("ent is not used in the project")
	}

	return util.NewTaskComposer("").AddFile(
		[]*task.FileDesc{
			{
				Path: path.Join("migrate", "config", "config.go"),
				Source: func(buf *bytes.Buffer) error {
					templates.WriteGoMigrateConfig(buf)
					return nil
				},
			},
			{
				Path: path.Join("migrate", "migrate.go"),
				Source: func(buf *bytes.Buffer) error {
					templates.WriteGoMigrateMigration(buf, meta.ProjectName, meta.Author)
					return nil
				},
				Transforms: []task.Pipeline{transformation.GoFormatSource},
			},
			{
				Path: path.Join("migrate", "wire.go"),
				Source: func(buf *bytes.Buffer) error {
					templates.WriteGoMigrateWire(buf, meta.ProjectName, meta.Author)
					return nil
				},
				Transforms: []task.Pipeline{transformation.GoFormatSource},
			},
			{
				Path: path.Join("cmd", "migrate", "main.go"),
				Source: func(buf *bytes.Buffer) error {
					templates.WriteGoMigrateCmdMain(buf, meta.ProjectName, meta.Author)
					return nil
				},
				Transforms: []task.Pipeline{transformation.GoFormatSource},
			},
			{
				Path: path.Join("pkg", "db", "ent.go"),
				Source: func(buf *bytes.Buffer) error {
					templates.WriteGoMigratePkgEnt(buf, meta.ProjectName, meta.Author)
					return nil
				},
				Transforms: []task.Pipeline{transformation.GoFormatSource},
			},
		}...).AddCommand(
		commands.GoModTidy(),
		commands.GoWire("./migrate"),
	).Run()
})

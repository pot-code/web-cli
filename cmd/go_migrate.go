package cmd

import (
	"bytes"

	"github.com/pkg/errors"
	"github.com/pot-code/web-cli/internal/command"
	"github.com/pot-code/web-cli/internal/constant"
	"github.com/pot-code/web-cli/internal/shell"
	"github.com/pot-code/web-cli/internal/task"
	"github.com/pot-code/web-cli/internal/transformer"
	"github.com/pot-code/web-cli/internal/util"
	"github.com/pot-code/web-cli/templates"
	"github.com/urfave/cli/v2"
)

type GoMigrateConfig struct{}

var GoMigrateCmd = command.NewCliCommand("migrate", "add migration",
	&GoMigrateConfig{},
	command.WithAlias([]string{"M"}),
).AddFeature(AddGoMigration).ExportCommand()

var AddGoMigration = command.NoCondFeature(func(c *cli.Context, cfg interface{}) error {
	meta, err := util.ParseGoMod(constant.GoModFile)
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

	return task.NewSequentialExecutor(
		[]task.Task{
			task.NewParallelExecutor(
				task.BatchFileGenerationTask(
					task.NewFileGenerationTree("").
						Branch("migrate").
						AddNodes( // migrate
							&task.FileGenerator{
								Name: "migrate.go",
								Data: bytes.NewBufferString(templates.GoMigrateMigration(meta.ProjectName, meta.Author)),
							},
							&task.FileGenerator{
								Name: "wire.go",
								Data: bytes.NewBufferString(templates.GoMigrateWire(meta.ProjectName, meta.Author)),
							},
						).
						Branch("config").
						AddNodes( // migrate/config
							&task.FileGenerator{
								Name: "config.go",
								Data: bytes.NewBufferString(templates.GoMigrateConfig()),
							},
						).Up().Up().
						Branch("cmd").
						AddNodes( // cmd
							&task.FileGenerator{
								Name: "main.go",
								Data: bytes.NewBufferString(templates.GoMigrateCmdMain(meta.ProjectName, meta.Author)),
							},
						).Up().
						Branch("pkg").Branch("db").
						AddNodes( // pkg/db
							&task.FileGenerator{
								Name: "ent.go",
								Data: bytes.NewBufferString(templates.GoMigratePkgEnt(meta.ProjectName, meta.Author)),
							},
						).
						Flatten(),
					transformer.GoFormatSource(),
				),
			),
			shell.GoModTidy(),
			shell.GoWire("./migrate"),
		},
	).Run()
})

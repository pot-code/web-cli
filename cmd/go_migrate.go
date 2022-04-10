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

var GoMigrateCmd = command.NewCliCommand("migrate", "add a Ent migration module",
	&GoMigrateConfig{},
	command.WithAlias([]string{"M"}),
).AddHandlers(AddGoMigration).BuildCommand()

var AddGoMigration = command.InlineHandler(func(c *cli.Context, cfg interface{}) error {
	meta, err := util.ParseGoMod(constant.GoModFile)
	if err != nil {
		return errors.WithStack(err)
	}

	if !meta.Contains("entgo.io/ent") {
		return errors.New("ent is not used in the project")
	}

	author := meta.GetAuthor()
	project := meta.GetProject()

	return task.NewSequentialExecutor(
		[]task.Task{
			task.NewParallelExecutor(
				task.BatchFileGenerationTask(
					task.NewFileGenerationTree("").
						Branch("migrate").
						AddNodes( // migrate
							&task.FileGenerator{
								Name: "migrate.go",
								Data: bytes.NewBufferString(templates.GoMigrateMigration(project, author)),
							},
							&task.FileGenerator{
								Name: "wire.go",
								Data: bytes.NewBufferString(templates.GoMigrateWire(project, author)),
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
								Data: bytes.NewBufferString(templates.GoMigrateCmdMain(project, author)),
							},
						).Up().
						Branch("pkg").Branch("db").
						AddNodes( // pkg/db
							&task.FileGenerator{
								Name: "ent.go",
								Data: bytes.NewBufferString(templates.GoMigratePkgEnt(project, author)),
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

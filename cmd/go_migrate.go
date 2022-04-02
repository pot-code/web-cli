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

var AddGoMigration = util.NoCondFeature(func(c *cli.Context, cfg interface{}) error {
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
		task.NewParallelExecutor(
			task.BatchFileTask(
				task.NewFileRequestTree("").
					Branch("migrate").Branch("config").
					AddNode( // migrate/config
						&task.FileRequest{
							Name: "config.go",
							Data: bytes.NewBufferString(templates.GoMigrateConfig()),
						},
					).Up().
					AddNode( // migrate/
						&task.FileRequest{
							Name: "migrate.go",
							Data: bytes.NewBufferString(templates.GoMigrateMigration(meta.ProjectName, meta.Author)),
						},
						&task.FileRequest{
							Name: "wire.go",
							Data: bytes.NewBufferString(templates.GoMigrateWire(meta.ProjectName, meta.Author)),
						},
					).Up().
					Branch("cmd").
					AddNode( // cmd
						&task.FileRequest{
							Name: "main.go",
							Data: bytes.NewBufferString(templates.GoMigrateCmdMain(meta.ProjectName, meta.Author)),
						},
					).Up().
					Branch("pkg").Branch("db").
					AddNode( // pkg/db
						&task.FileRequest{
							Name: "ent.go",
							Data: bytes.NewBufferString(templates.GoMigratePkgEnt(meta.ProjectName, meta.Author)),
						}).
					Flatten(),
				transformer.GoFormatSource(),
			)...,
		),
		task.NewShellCmdExecutor(
			shell.GoModTidy(),
		),
		task.NewShellCmdExecutor(
			shell.GoWire("./migrate"),
		),
	).Run()
})

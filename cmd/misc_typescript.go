package cmd

import (
	"bytes"

	"github.com/pot-code/web-cli/internal/command"
	"github.com/pot-code/web-cli/internal/pkm"
	"github.com/pot-code/web-cli/internal/task"
	"github.com/pot-code/web-cli/templates"
	"github.com/urfave/cli/v2"
)

type AddTypescriptConfig struct {
	Target         string `flag:"target" alias:"t" usage:"project type" validate:"required,oneof=node react"`
	PackageManager string `flag:"pm" usage:"choose package manager" validate:"oneof=pnpm npm yarn"`
}

var AddTypescriptCmd = command.NewCliCommand("typescript", "add typescript support",
	&AddTypescriptConfig{
		PackageManager: "pnpm",
	},
	command.WithAlias([]string{"ts"}),
).AddHandlers(
	new(AddTypescript),
).BuildCommand()

type AddTypescript struct {
	PackageManager string
}

var _ command.CommandHandler = &AddTypescript{}

func (at *AddTypescript) Handle(c *cli.Context, cfg interface{}) error {
	config := cfg.(*AddTypescriptConfig)

	at.PackageManager = config.PackageManager

	if config.Target == "node" {
		return at.node().Run()
	}
	if config.Target == "react" {
		return at.react().Run()
	}
	return nil
}

func (at *AddTypescript) node() task.Task {
	pm := pkm.NewPackageManager(at.PackageManager)

	return task.NewParallelExecutor(
		[]task.Task{
			&task.FileGenerator{
				Name: ".eslintrc.js",
				Data: bytes.NewBufferString(templates.NodeEslintrc()),
			},
			&task.FileGenerator{
				Name: "tsconfig.json",
				Data: bytes.NewBufferString(templates.NodeTsConfig()),
			},
			pm.InstallDev(
				[]string{
					"typescript",
					"eslint",
					"@typescript-eslint/eslint-plugin",
					"eslint-plugin-prettier",
					"@typescript-eslint/parser",
					"eslint-config-prettier",
					"eslint-plugin-import",
					"prettier",
					"prettier-eslint",
				},
			),
		},
	)
}

func (at *AddTypescript) react() task.Task {
	pm := pkm.NewPackageManager(at.PackageManager)

	return task.NewParallelExecutor(
		[]task.Task{
			&task.FileGenerator{
				Name: ".eslintrc.js",
				Data: bytes.NewBufferString(templates.ReactEslintrc()),
			},
			&task.FileGenerator{
				Name: "tsconfig.json",
				Data: bytes.NewBufferString(templates.ReactTsConfig()),
			},
			pm.InstallDev(
				[]string{
					"@types/react",
					"@typescript-eslint/eslint-plugin",
					"@typescript-eslint/parser",
					"eslint",
					"eslint-config-airbnb",
					"eslint-config-prettier",
					"eslint-import-resolver-typescript",
					"eslint-plugin-import",
					"eslint-plugin-jsx-a11y",
					"eslint-plugin-prettier",
					"eslint-plugin-react",
					"eslint-plugin-react-hooks",
					"prettier",
					"prettier-eslint",
					"typescript",
				},
			),
		},
	)
}

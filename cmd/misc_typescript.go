package cmd

import (
	"bytes"

	"github.com/pot-code/web-cli/internal/command"
	"github.com/pot-code/web-cli/internal/shell"
	"github.com/pot-code/web-cli/internal/task"
	"github.com/pot-code/web-cli/templates"
	"github.com/urfave/cli/v2"
)

type AddTypescriptConfig struct {
	Target string `flag:"target" alias:"t" usage:"project target" validate:"required,oneof=node react"`
}

var AddTypescriptCmd = command.NewCliCommand("typescript", "add typescript support",
	new(AddTypescriptConfig),
	command.WithAlias([]string{"ts"}),
).AddFeature(
	new(AddTypescriptToNode),
	new(AddTypescriptToReact),
).ExportCommand()

type AddTypescriptToNode struct{}

var _ command.CommandFeature = &AddTypescriptToNode{}

func (arc *AddTypescriptToNode) Cond(c *cli.Context) bool {
	return c.String("target") == "node"
}

func (arc *AddTypescriptToNode) Handle(c *cli.Context, cfg interface{}) error {
	return task.NewParallelExecutor(
		&task.FileGenerator{
			Name: ".eslintrc.js",
			Data: bytes.NewBufferString(templates.NodeEslintrc()),
		},
		&task.FileGenerator{
			Name: "tsconfig.json",
			Data: bytes.NewBufferString(templates.NodeTsConfig()),
		},
		shell.YarnAddDev(
			"typescript",
			"eslint",
			"@typescript-eslint/eslint-plugin",
			"eslint-plugin-prettier",
			"@typescript-eslint/parser",
			"eslint-config-prettier",
			"eslint-plugin-import",
			"prettier",
			"prettier-eslint",
		),
	).Run()
}

type AddTypescriptToReact struct{}

var _ command.CommandFeature = &AddTypescriptToReact{}

func (arc *AddTypescriptToReact) Cond(c *cli.Context) bool {
	return c.String("target") == "react"
}

func (arc *AddTypescriptToReact) Handle(c *cli.Context, cfg interface{}) error {
	return task.NewParallelExecutor(
		&task.FileGenerator{
			Name: ".eslintrc.js",
			Data: bytes.NewBufferString(templates.ReactEslintrc()),
		},
		&task.FileGenerator{
			Name: "tsconfig.json",
			Data: bytes.NewBufferString(templates.ReactTsConfig()),
		},
		shell.YarnAddDev(
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
		),
	).Run()
}

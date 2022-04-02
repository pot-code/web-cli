package cmd

import (
	"bytes"

	"github.com/pot-code/web-cli/internal/commands"
	"github.com/pot-code/web-cli/internal/task"
	"github.com/pot-code/web-cli/internal/util"
	"github.com/pot-code/web-cli/templates"
	"github.com/urfave/cli/v2"
)

type AddTypescriptConfig struct {
	Target string `flag:"target" alias:"t" usage:"project target" validate:"required,oneof=node react"`
}

var AddTypescriptCmd = util.NewCliCommand("typescript", "add typescript support",
	new(AddTypescriptConfig),
	util.WithAlias([]string{"ts"}),
).AddFeature(
	new(AddTypescriptToNode),
	new(AddTypescriptToReact),
).ExportCommand()

type AddTypescriptToNode struct{}

var _ util.CommandFeature = &AddTypescriptToNode{}

func (arc *AddTypescriptToNode) Cond(c *cli.Context) bool {
	return c.String("target") == "node"
}

func (arc *AddTypescriptToNode) Handle(c *cli.Context, cfg interface{}) error {
	return util.NewTaskComposer("").AddFile(
		&task.FileDesc{
			Path: ".eslintrc.js",
			Source: func(buf *bytes.Buffer) error {
				templates.WriteNodeEslintrc(buf)
				return nil
			},
		},
		&task.FileDesc{
			Path: "tsconfig.json",
			Source: func(buf *bytes.Buffer) error {
				templates.WriteNodeTsConfig(buf)
				return nil
			},
		},
	).AddCommand(
		commands.YarnAddDev(
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

var _ util.CommandFeature = &AddTypescriptToReact{}

func (arc *AddTypescriptToReact) Cond(c *cli.Context) bool {
	return c.String("target") == "react"
}

func (arc *AddTypescriptToReact) Handle(c *cli.Context, cfg interface{}) error {
	return util.NewTaskComposer("").AddFile(
		&task.FileDesc{
			Path: ".eslintrc.js",
			Source: func(buf *bytes.Buffer) error {
				templates.WriteReactEslintrc(buf)
				return nil
			},
		},
		&task.FileDesc{
			Path: "tsconfig.json",
			Source: func(buf *bytes.Buffer) error {
				templates.WriteReactTsConfig(buf)
				return nil
			},
		},
	).AddCommand(
		commands.YarnAddDev(
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

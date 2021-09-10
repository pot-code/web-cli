package cmd

import (
	"bytes"

	"github.com/pot-code/web-cli/pkg/core"
	"github.com/pot-code/web-cli/pkg/util"
	"github.com/pot-code/web-cli/templates"
	"github.com/urfave/cli/v2"
)

type AddTypescriptConfig struct {
	Target string `flag:"target" alias:"t" usage:"project target" validate:"required,oneof=node react"`
}

var AddTypescriptCmd = core.NewCliLeafCommand("typescript", "add typescript support", new(AddTypescriptConfig),
	core.WithAlias([]string{"ts"}),
).AddService(new(AddTypescriptToNodeService)).AddService(new(AddTypescriptToReactService)).ExportCommand()

type AddTypescriptToNodeService struct{}

var _ core.CommandService = &AddTypescriptToNodeService{}

func (arc *AddTypescriptToNodeService) Cond(c *cli.Context) bool {
	return c.String("target") == "node"
}

func (arc *AddTypescriptToNodeService) Handle(c *cli.Context, cfg interface{}) error {
	return util.NewTaskComposer("").AddFile(
		&core.FileDesc{
			Path: ".eslintrc.js",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteNodeEslintrc(&buf)
				return buf.Bytes()
			},
		},
		&core.FileDesc{
			Path: "tsconfig.json",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteNodeTsConfig(&buf)
				return buf.Bytes()
			},
		},
	).AddCommand(
		&core.Command{
			Bin: "npm",
			Args: []string{"i", "-D",
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
		},
	).Run()
}

type AddTypescriptToReactService struct{}

var _ core.CommandService = &AddTypescriptToReactService{}

func (arc *AddTypescriptToReactService) Cond(c *cli.Context) bool {
	return c.String("target") == "react"
}

func (arc *AddTypescriptToReactService) Handle(c *cli.Context, cfg interface{}) error {
	return util.NewTaskComposer("").AddFile(
		&core.FileDesc{
			Path: ".eslintrc.js",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteReactEslintrc(&buf)
				return buf.Bytes()
			},
		},
		&core.FileDesc{
			Path: "tsconfig.json",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteReactTsConfig(&buf)
				return buf.Bytes()
			},
		},
	).AddCommand(
		&core.Command{
			Bin: "npm",
			Args: []string{"i", "-D",
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
		},
	).Run()
}

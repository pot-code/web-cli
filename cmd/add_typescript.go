package cmd

import (
	"bytes"

	"github.com/pot-code/web-cli/pkg/core"
	"github.com/pot-code/web-cli/pkg/util"
	"github.com/pot-code/web-cli/templates"
	"github.com/urfave/cli/v2"
)

const cmdTypescriptName = "typescript"

type addTypescriptConfig struct {
	Target string `name:"target" validate:"required,oneof=node react"`
}

var addTypescriptCmd = &cli.Command{
	Name:    cmdTypescriptName,
	Usage:   "add typescript support",
	Aliases: []string{"ts"},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "target",
			Aliases: []string{"t"},
			Usage:   "project target (node/react)",
			Value:   "node",
		},
	},
	Action: func(c *cli.Context) error {
		config := new(addTypescriptConfig)
		err := util.ParseConfig(c, config)
		if err != nil {
			if _, ok := err.(*util.CommandError); ok {
				cli.ShowCommandHelp(c, cmdTypescriptName)
			}
			return err
		}

		var cmd core.Generator
		if config.Target == "node" {
			cmd = newAddTypescriptToNode()
		} else {
			cmd = newAddTypescriptToReact()
		}

		err = cmd.Run()
		if err != nil {
			cmd.Cleanup()
		}
		return err
	},
}

func newAddTypescriptToNode() core.Generator {
	return util.NewTaskComposer("",
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
	).AddCommand(&core.Command{
		Bin: "npm",
		Args: []string{"i", "-D",
			"typescript",
			"eslint",
			"@typescript-eslint/eslint-plugin",
			"eslint-plugin-prettier",
			"@typescript-eslint/parser",
			"eslint-config-prettier",
			"eslint-plugin-import",
		},
	})
}

func newAddTypescriptToReact() core.Generator {
	return util.NewTaskComposer("",
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
	).AddCommand(&core.Command{
		Bin: "npm",
		Args: []string{"i", "-D",
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
	})
}

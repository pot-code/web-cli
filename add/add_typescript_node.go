package add

import (
	"bytes"

	"github.com/pkg/errors"
	"github.com/pot-code/web-cli/core"
	"github.com/pot-code/web-cli/templates"
	"github.com/pot-code/web-cli/util"
)

type AddTypescriptToNode struct {
	runner core.Generator
}

var _ core.Generator = AddTypescriptToNode{}

func NewAddTypescriptToNode() *AddTypescriptToNode {
	composer := util.NewTaskComposer("",
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
	)
	composer.AddCommand(&core.Command{
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
	return &AddTypescriptToNode{runner: composer}
}

func (atn AddTypescriptToNode) Run() error {
	return errors.Wrap(atn.runner.Run(), "failed to generate typescript config")
}

func (atn AddTypescriptToNode) Cleanup() error {
	return atn.runner.Cleanup()
}

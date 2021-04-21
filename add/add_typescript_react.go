package add

import (
	"bytes"

	"github.com/pkg/errors"
	"github.com/pot-code/web-cli/core"
	"github.com/pot-code/web-cli/templates"
	"github.com/pot-code/web-cli/util"
)

type AddTypescriptToReact struct {
	runner core.Generator
}

var _ core.Generator = AddTypescriptToReact{}

func NewAddTypescriptToReact() *AddTypescriptToReact {
	composer := util.NewTaskComposer("",
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
	)
	composer.AddCommand(&core.Command{
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
	return &AddTypescriptToReact{runner: composer}
}

func (atr AddTypescriptToReact) Run() error {
	return errors.Wrap(atr.runner.Run(), "failed to generate typescript config")
}

func (atr AddTypescriptToReact) Cleanup() error {
	return atr.runner.Cleanup()
}

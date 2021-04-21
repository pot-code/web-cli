package add

import (
	"bytes"

	"github.com/pkg/errors"
	"github.com/pot-code/web-cli/core"
	"github.com/pot-code/web-cli/templates"
	"github.com/pot-code/web-cli/util"
	log "github.com/sirupsen/logrus"
)

type AddTypescriptToReact struct {
	composer *util.TaskComposer
	runner   core.Generator
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
	return &AddTypescriptToReact{composer: composer}
}

func (atr AddTypescriptToReact) Run() error {
	log.Debugf("runnereration tree:\n%s", atr.composer.GetGenerationTree())
	runner := atr.composer.MakeRunner()
	atr.runner = runner
	return errors.Wrap(runner.Run(), "failed to runnererate typescript config")
}

func (atr AddTypescriptToReact) Cleanup() error {
	return atr.runner.Cleanup()
}

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
	recipe *util.GenerationRecipe
	gen    core.Generator
}

var _ core.Executor = AddTypescriptToReact{}

func NewAddTypescriptToReact() *AddTypescriptToReact {
	recipe := util.NewGenerationRecipe("",
		&util.GenerationMaterial{
			Path: "./.eslintrc.js",
			Provider: func() []byte {
				var buf bytes.Buffer

				templates.WriteReactEslintrc(&buf)
				return buf.Bytes()
			},
		},
		&util.GenerationMaterial{
			Path: "./tsconfig.json",
			Provider: func() []byte {
				var buf bytes.Buffer

				templates.WriteReactTsConfig(&buf)
				return buf.Bytes()
			},
		},
	)
	return &AddTypescriptToReact{recipe: recipe}
}

func (atn AddTypescriptToReact) Run() error {
	cmd := core.NewCmdExecutor("npm", "i", "-D",
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
	)

	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "failed to install dependencies")
	}

	log.Debugf("generation tree:\n%s", atn.recipe.GetGenerationTree())
	gen := atn.recipe.MakeGenerator()
	atn.gen = gen
	return errors.Wrap(gen.Run(), "failed to generate typescript config")
}

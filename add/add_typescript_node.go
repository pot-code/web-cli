package add

import (
	"bytes"

	"github.com/pkg/errors"
	"github.com/pot-code/web-cli/core"
	"github.com/pot-code/web-cli/templates"
	"github.com/pot-code/web-cli/util"
	log "github.com/sirupsen/logrus"
)

type AddTypescriptToNode struct {
	recipe *util.GenerationRecipe
	gen    core.Generator
}

var _ core.Executor = AddTypescriptToNode{}

func NewAddTypescriptToNode() *AddTypescriptToNode {
	recipe := util.NewGenerationRecipe("",
		&util.GenerationMaterial{
			Path: ".eslintrc.js",
			Provider: func() []byte {
				var buf bytes.Buffer

				templates.WriteNodeEslintrc(&buf)
				return buf.Bytes()
			},
		},
		&util.GenerationMaterial{
			Path: "tsconfig.json",
			Provider: func() []byte {
				var buf bytes.Buffer

				templates.WriteNodeTsConfig(&buf)
				return buf.Bytes()
			},
		},
	)
	return &AddTypescriptToNode{recipe: recipe}
}

func (atn AddTypescriptToNode) Run() error {
	cmd := core.NewCmdExecutor("npm", "i", "-D",
		"typescript",
		"eslint",
		"@typescript-eslint/eslint-plugin",
		"eslint-plugin-prettier",
		"@typescript-eslint/parser",
		"eslint-config-prettier",
		"eslint-plugin-import",
	)

	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "failed to install dependencies")
	}

	log.Debugf("generation tree:\n%s", atn.recipe.GetGenerationTree())
	gen := atn.recipe.MakeGenerator()
	atn.gen = gen
	return errors.Wrap(gen.Run(), "failed to generate typescript config")
}

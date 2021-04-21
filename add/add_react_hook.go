package add

import (
	"bytes"
	"fmt"

	"github.com/pkg/errors"
	"github.com/pot-code/web-cli/core"
	"github.com/pot-code/web-cli/templates"
	"github.com/pot-code/web-cli/util"
	log "github.com/sirupsen/logrus"
)

type AddReactHook struct {
	recipe *util.GenerationRecipe
	gen    core.Generator
}

var _ core.Executor = AddReactHook{}

func NewAddReactHook(name string) *AddReactHook {
	recipe := util.NewGenerationRecipe("",
		&util.GenerationMaterial{
			Path: fmt.Sprintf("./%s.%s", name, "ts"),
			Provider: func() []byte {
				var buf bytes.Buffer

				templates.WriteReactHook(&buf, name)
				return buf.Bytes()
			},
		},
	)
	return &AddReactHook{recipe: recipe}
}

func (atn AddReactHook) Run() error {
	log.Debugf("generation tree:\n%s", atn.recipe.GetGenerationTree())
	gen := atn.recipe.MakeGenerator()
	atn.gen = gen
	return errors.Wrap(gen.Run(), "failed to generate react hook")
}

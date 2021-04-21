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

type AddReactComponent struct {
	recipe *util.GenerationRecipe
	gen    core.Generator
}

var _ core.Executor = AddReactComponent{}

func NewAddReactComponent(name string) *AddReactComponent {
	recipe := util.NewGenerationRecipe("",
		&util.GenerationMaterial{
			Path: fmt.Sprintf("./%s.%s", name, "tsx"),
			Provider: func() []byte {
				var buf bytes.Buffer

				templates.WriteReactComponent(&buf, name)
				return buf.Bytes()
			},
		},
	)
	return &AddReactComponent{recipe: recipe}
}

func (atn AddReactComponent) Run() error {
	log.Debugf("generation tree:\n%s", atn.recipe.GetGenerationTree())
	gen := atn.recipe.MakeGenerator()
	atn.gen = gen
	return errors.Wrap(gen.Run(), "failed to generate react component")
}

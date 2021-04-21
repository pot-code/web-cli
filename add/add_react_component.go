package add

import (
	"bytes"
	"fmt"

	"github.com/pkg/errors"
	"github.com/pot-code/web-cli/core"
	"github.com/pot-code/web-cli/templates"
	"github.com/pot-code/web-cli/util"
)

type AddReactComponent struct {
	runner core.Generator
}

var _ core.Executor = AddReactComponent{}

func NewAddReactComponent(name string) *AddReactComponent {
	composer := util.NewTaskComposer("",
		&core.FileDesc{
			Path: fmt.Sprintf("%s.%s", name, "tsx"),
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteReactComponent(&buf, name)
				return buf.Bytes()
			},
		},
	)
	return &AddReactComponent{runner: composer}
}

func (arc AddReactComponent) Run() error {
	return errors.Wrap(arc.runner.Run(), "failed to generate react component")
}

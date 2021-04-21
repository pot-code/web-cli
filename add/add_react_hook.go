package add

import (
	"bytes"
	"fmt"

	"github.com/pkg/errors"
	"github.com/pot-code/web-cli/core"
	"github.com/pot-code/web-cli/templates"
	"github.com/pot-code/web-cli/util"
)

type AddReactHook struct {
	runner core.Generator
}

var _ core.Executor = AddReactHook{}

func NewAddReactHook(name string) *AddReactHook {
	composer := util.NewTaskComposer("",
		&core.FileDesc{
			Path: fmt.Sprintf("%s.%s", name, "ts"),
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteReactHook(&buf, name)
				return buf.Bytes()
			},
		},
	)
	return &AddReactHook{runner: composer}
}

func (arh AddReactHook) Run() error {
	return errors.Wrap(arh.runner.Run(), "failed to generate react hook")
}

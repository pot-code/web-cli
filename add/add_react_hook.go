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
	composer *util.TaskComposer
	runner   core.Generator
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
	return &AddReactHook{composer: composer}
}

func (atn AddReactHook) Run() error {
	log.Debugf("runnereration tree:\n%s", atn.composer.GetGenerationTree())
	runner := atn.composer.MakeRunner()
	atn.runner = runner
	return errors.Wrap(runner.Run(), "failed to runnererate react hook")
}

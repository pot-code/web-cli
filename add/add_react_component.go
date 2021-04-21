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
	composer *util.TaskComposer
	runner   core.Generator
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
	return &AddReactComponent{composer: composer}
}

func (atn AddReactComponent) Run() error {
	log.Debugf("runnereration tree:\n%s", atn.composer.GetGenerationTree())
	runner := atn.composer.MakeRunner()
	atn.runner = runner
	return errors.Wrap(runner.Run(), "failed to runnererate react component")
}

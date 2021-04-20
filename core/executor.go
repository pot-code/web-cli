package core

import (
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type CmdExecutor struct {
	name string
	args []string
}

var _ Executor = CmdExecutor{}

func NewCmdExecutor(bin string, args ...string) *CmdExecutor {
	return &CmdExecutor{bin, args}
}

func (ce CmdExecutor) Run() error {
	log.WithFields(log.Fields{
		"bin":  ce.name,
		"args": ce.args,
	}).Debug("execute command")
	cmd := exec.Command(ce.name, ce.args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	return errors.Wrapf(cmd.Run(), "failed to execute command '%s %s'", ce.name, strings.Join(ce.args, " "))
}

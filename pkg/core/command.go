package core

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Command struct {
	Bin  string
	Dir  string
	Args []string
}

func (c Command) String() string {
	if len(c.Args) > 0 {
		return fmt.Sprintf("%s %s", c.Bin, strings.Join(c.Args, " "))
	}
	return c.Bin
}

type CmdExecutor struct {
	cmd *Command
}

var _ Executor = CmdExecutor{}

func NewCmdExecutor(cmd *Command) *CmdExecutor {
	return &CmdExecutor{cmd}
}

func (ce CmdExecutor) Run() error {
	log.WithField("cmd", ce.cmd).Info("execute command")
	cmd := ce.cmd
	proc := exec.Command(cmd.Bin, cmd.Args...)

	if ce.cmd.Dir != "" {
		proc.Dir = ce.cmd.Dir
	}
	proc.Stdout = os.Stdout
	proc.Stderr = os.Stdout
	return errors.Wrapf(proc.Run(), "failed to execute command '%s'", cmd)
}

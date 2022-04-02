package task

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type ShellCommand struct {
	Bin  string
	Cwd  string
	Args []string
	Out  io.Writer
}

func (c ShellCommand) String() string {
	return fmt.Sprintf("cwd=%s bin=%s args=%s", c.Cwd, c.Bin, strings.Join(c.Args, " "))
}

type ShellCmdExecutor struct {
	cmd *ShellCommand
}

var _ Task = &ShellCmdExecutor{}

func NewShellCmdExecutor(cmd *ShellCommand) *ShellCmdExecutor {
	return &ShellCmdExecutor{cmd}
}

func (ce *ShellCmdExecutor) Run() error {
	cmd := ce.cmd
	proc := exec.Command(cmd.Bin, cmd.Args...)

	log.WithField("cmd", cmd).Info("execute command")

	if cmd.Cwd != "" {
		proc.Dir = cmd.Cwd
	}

	proc.Stdout = os.Stdout
	if cmd.Out != nil {
		proc.Stdout = cmd.Out
	}
	proc.Stderr = os.Stderr
	return errors.Wrapf(proc.Run(), "failed to execute command '%s'", cmd)
}

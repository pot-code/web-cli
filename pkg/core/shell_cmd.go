package core

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
	Bin    string
	Cwd    string
	Before bool // run before file generation
	Args   []string
	Out    io.Writer
}

func (c *ShellCommand) String() string {
	return fmt.Sprintf("cwd=%s bin=%s args=%s", c.Cwd, c.Bin, strings.Join(c.Args, ","))
}

type ShellCmdExecutor struct {
	cmd *ShellCommand
}

var _ Runner = ShellCmdExecutor{}

func NewShellCmdExecutor(cmd *ShellCommand) *ShellCmdExecutor {
	return &ShellCmdExecutor{cmd}
}

func (ce ShellCmdExecutor) Run() error {
	cmd := ce.cmd
	log.WithFields(log.Fields{"cmd": cmd, "context": "ShellCmdExecutor.Run"}).Info("execute command")
	proc := exec.Command(cmd.Bin, cmd.Args...)

	if cmd.Cwd != "" {
		proc.Dir = cmd.Cwd
	}

	proc.Stdout = os.Stdout
	if cmd.Out != nil {
		proc.Stdout = cmd.Out
	}
	proc.Stderr = os.Stdout
	return errors.Wrapf(proc.Run(), "failed to execute command '%s'", cmd)
}

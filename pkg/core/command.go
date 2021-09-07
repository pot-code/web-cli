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

type Command struct {
	Bin    string
	Dir    string
	Before bool // run before file generation
	Args   []string
	Out    io.Writer
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

var _ Runner = CmdExecutor{}

func NewCmdExecutor(cmd *Command) *CmdExecutor {
	return &CmdExecutor{cmd}
}

func (ce CmdExecutor) Run() error {
	cmd := ce.cmd
	log.WithFields(log.Fields{"cmd": cmd, "context": "CmdExecutor.Run"}).Info("execute command")
	proc := exec.Command(cmd.Bin, cmd.Args...)

	if cmd.Dir != "" {
		proc.Dir = cmd.Dir
	}

	proc.Stdout = os.Stdout
	if cmd.Out != nil {
		proc.Stdout = cmd.Out
	}
	proc.Stderr = os.Stdout
	return errors.Wrapf(proc.Run(), "failed to execute command '%s'", cmd)
}

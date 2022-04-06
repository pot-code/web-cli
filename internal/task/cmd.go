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

func (c *ShellCommand) String() string {
	return fmt.Sprintf("cwd=%s bin=%s args=%s", c.Cwd, c.Bin, strings.Join(c.Args, " "))
}

var _ Task = &ShellCommand{}

func (sc *ShellCommand) Run() error {
	proc := exec.Command(sc.Bin, sc.Args...)

	log.WithField("cmd", sc).Info("execute command")

	if sc.Cwd != "" {
		proc.Dir = sc.Cwd
	}

	proc.Stdout = os.Stdout
	if sc.Out != nil {
		proc.Stdout = sc.Out
	}
	proc.Stderr = os.Stderr
	return errors.Wrapf(proc.Run(), "failed to execute command '%s'", sc)
}
package task

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

type ShellCommandTask struct {
	Bin  string
	Cwd  string
	Args []string
	Out  io.Writer
}

func (c *ShellCommandTask) String() string {
	return fmt.Sprintf("%s %s", c.Bin, strings.Join(c.Args, " "))
}

var _ Task = (*ShellCommandTask)(nil)

func (sc *ShellCommandTask) Run() error {
	proc := exec.Command(sc.Bin, sc.Args...)
	log.WithField("cwd", sc.Cwd).Infof("run shell command '%s'", sc)

	if sc.Cwd != "" {
		proc.Dir = sc.Cwd
	}

	proc.Stdout = os.Stdout
	if sc.Out != nil {
		proc.Stdout = sc.Out
	}
	proc.Stderr = os.Stderr

	err := proc.Run()
	if err != nil {
		return fmt.Errorf("run command [cmd: %s]: %w", sc, err)
	}
	return nil
}

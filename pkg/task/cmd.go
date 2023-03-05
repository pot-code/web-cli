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
	bin  string
	cwd  string
	args []string
	out  io.Writer
}

func (t *ShellCommandTask) String() string {
	return fmt.Sprintf("%s %s", t.bin, strings.Join(t.args, " "))
}

var _ Task = (*ShellCommandTask)(nil)

func (t *ShellCommandTask) Run() error {
	proc := exec.Command(t.bin, t.args...)
	log.WithField("cwd", t.cwd).Infof("run shell command '%s'", t)

	if t.cwd != "" {
		proc.Dir = t.cwd
	}

	proc.Stdout = os.Stdout
	if t.out != nil {
		proc.Stdout = t.out
	}
	proc.Stderr = os.Stderr

	err := proc.Run()
	if err != nil {
		return fmt.Errorf("run command [cmd: %s]: %w", t, err)
	}
	return nil
}

package task

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
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
	log.Info().Str("task", "ShellCommandTask").Str("cwd", t.cwd).Str("cmd", t.String()).Msg("run shell command")

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

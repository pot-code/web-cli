package commands

import (
	"path"

	"github.com/pot-code/web-cli/pkg/task"
)

func GitClone(url, dst string) *task.ShellCommand {
	return &task.ShellCommand{
		Bin:  "git",
		Args: []string{"clone", url, dst},
	}
}

func GitDeleteHistory(dst string) *task.ShellCommand {
	return &task.ShellCommand{
		Bin:  "rm",
		Args: []string{"-rf", path.Join(dst, ".git")},
	}
}

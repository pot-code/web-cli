package shell

import (
	"path"

	"github.com/pot-code/web-cli/internal/task"
)

func GitClone(url, dst string) *task.ShellCommandTask {
	return &task.ShellCommandTask{
		Bin:  "git",
		Args: []string{"clone", url, dst},
	}
}

func GitDeleteHistory(dst string) *task.ShellCommandTask {
	return &task.ShellCommandTask{
		Bin:  "rm",
		Args: []string{"-rf", path.Join(dst, ".git")},
	}
}

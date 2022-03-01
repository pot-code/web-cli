package commands

import (
	"path"

	"github.com/pot-code/web-cli/pkg/core"
)

func GitClone(url, dst string) *core.ShellCommand {
	return &core.ShellCommand{
		Bin:  "git",
		Args: []string{"clone", url, dst},
	}
}

func GitDeleteHistory(dst string) *core.ShellCommand {
	return &core.ShellCommand{
		Bin:  "rm",
		Args: []string{"-rf", path.Join(dst, ".git")},
	}
}

package commands

import "github.com/pot-code/web-cli/pkg/core"

func GitClone(url, dest string) *core.ShellCommand {
	return &core.ShellCommand{
		Bin:  "git",
		Args: []string{"clone", url, dest},
	}
}

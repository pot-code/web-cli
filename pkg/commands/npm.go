package commands

import "github.com/pot-code/web-cli/pkg/core"

func NpmAdd(deps ...string) *core.ShellCommand {
	args := []string{"i"}
	args = append(args, deps...)
	return &core.ShellCommand{
		Bin:  "npm",
		Args: args,
	}
}

func NpmAddDev(deps ...string) *core.ShellCommand {
	args := []string{"i", "-D"}
	args = append(args, deps...)
	return &core.ShellCommand{
		Bin:  "npm",
		Args: args,
	}
}

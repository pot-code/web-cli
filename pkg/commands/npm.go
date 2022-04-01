package commands

import "github.com/pot-code/web-cli/pkg/task"

func NpmAdd(deps ...string) *task.ShellCommand {
	args := []string{"i"}
	args = append(args, deps...)
	return &task.ShellCommand{
		Bin:  "npm",
		Args: args,
	}
}

func NpmAddDev(deps ...string) *task.ShellCommand {
	args := []string{"i", "-D"}
	args = append(args, deps...)
	return &task.ShellCommand{
		Bin:  "npm",
		Args: args,
	}
}

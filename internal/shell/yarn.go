package shell

import "github.com/pot-code/web-cli/internal/task"

func YarnAdd(deps ...string) *task.ShellCommand {
	args := []string{"add"}
	args = append(args, deps...)
	return &task.ShellCommand{
		Bin:  "yarn",
		Args: args,
	}
}

func YarnAddDev(deps ...string) *task.ShellCommand {
	args := []string{"add", "-D"}
	args = append(args, deps...)
	return &task.ShellCommand{
		Bin:  "yarn",
		Args: args,
	}
}

func YarnCreate(template string, extra ...string) *task.ShellCommand {
	args := []string{"create", template}
	args = append(args, extra...)
	return &task.ShellCommand{
		Bin:  "yarn",
		Args: args,
	}
}

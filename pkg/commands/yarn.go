package commands

import "github.com/pot-code/web-cli/pkg/core"

func YarnAdd(deps ...string) *core.ShellCommand {
	args := []string{"add"}
	args = append(args, deps...)
	return &core.ShellCommand{
		Bin:  "yarn",
		Args: args,
	}
}

func YarnAddDev(deps ...string) *core.ShellCommand {
	args := []string{"add", "-D"}
	args = append(args, deps...)
	return &core.ShellCommand{
		Bin:  "yarn",
		Args: args,
	}
}

func YarnCreate(template string, extra ...string) *core.ShellCommand {
	args := []string{"create", template}
	args = append(args, extra...)
	return &core.ShellCommand{
		Bin:  "yarn",
		Args: args,
	}
}

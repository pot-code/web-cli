package pkm

import "github.com/pot-code/web-cli/internal/task"

type yarn struct {
	bin string
}

func newYarn(bin string) *yarn {
	return &yarn{bin: bin}
}

func (y *yarn) Create(template, name string, flags []string) *task.ShellCommand {
	args := []string{"create", template}
	args = append(args, flags...)
	args = append(args, name)
	return &task.ShellCommand{
		Bin:  y.bin,
		Args: args,
	}
}

func (y *yarn) Install(name []string) *task.ShellCommand {
	args := []string{"add"}
	args = append(args, name...)
	return &task.ShellCommand{
		Bin:  y.bin,
		Args: args,
	}
}

func (y *yarn) InstallDev(name []string) *task.ShellCommand {
	args := []string{"add", "-D"}
	args = append(args, name...)
	return &task.ShellCommand{
		Bin:  y.bin,
		Args: args,
	}
}

func (y *yarn) Uninstall(name []string) *task.ShellCommand {
	args := []string{"remove"}
	args = append(args, name...)
	return &task.ShellCommand{
		Bin:  y.bin,
		Args: args,
	}
}

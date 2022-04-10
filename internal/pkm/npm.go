package pkm

import "github.com/pot-code/web-cli/internal/task"

type npm struct {
	bin string
}

func newNpm(bin string) *npm {
	return &npm{bin: bin}
}

func (n *npm) Create(template, name string, flags []string) *task.ShellCommand {
	bin := "npx"
	args := []string{template}
	args = append(args, flags...)
	args = append(args, name)
	return &task.ShellCommand{
		Bin:  bin,
		Args: args,
	}
}

func (n *npm) Install(name []string) *task.ShellCommand {
	args := []string{"install"}
	args = append(args, name...)
	return &task.ShellCommand{
		Bin:  n.bin,
		Args: args,
	}
}

func (n *npm) InstallDev(name []string) *task.ShellCommand {
	args := []string{"install", "-D"}
	args = append(args, name...)
	return &task.ShellCommand{
		Bin:  n.bin,
		Args: args,
	}
}

func (n *npm) Uninstall(name []string) *task.ShellCommand {
	args := []string{"rm"}
	args = append(args, name...)
	return &task.ShellCommand{
		Bin:  n.bin,
		Args: args,
	}
}

package pkm

import "github.com/pot-code/web-cli/internal/task"

type npm struct {
	bin string
}

func newNpm(bin string) *npm {
	return &npm{bin: bin}
}

func (n *npm) Create(template, name string, flags []string) *task.ShellCommandTask {
	bin := "npx"
	args := []string{template}
	args = append(args, flags...)
	args = append(args, name)
	return &task.ShellCommandTask{
		Bin:  bin,
		Args: args,
	}
}

func (n *npm) Install(name []string) *task.ShellCommandTask {
	args := []string{"install"}
	args = append(args, name...)
	return &task.ShellCommandTask{
		Bin:  n.bin,
		Args: args,
	}
}

func (n *npm) InstallDev(name []string) *task.ShellCommandTask {
	args := []string{"install", "-D"}
	args = append(args, name...)
	return &task.ShellCommandTask{
		Bin:  n.bin,
		Args: args,
	}
}

func (n *npm) Uninstall(name []string) *task.ShellCommandTask {
	args := []string{"rm"}
	args = append(args, name...)
	return &task.ShellCommandTask{
		Bin:  n.bin,
		Args: args,
	}
}

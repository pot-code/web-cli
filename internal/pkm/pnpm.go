package pkm

import "github.com/pot-code/web-cli/internal/task"

type pnpm struct {
	bin string
}

func newPnpm(bin string) *pnpm {
	return &pnpm{bin: bin}
}

func (p *pnpm) Create(template, name string, flags []string) *task.ShellCommand {
	args := []string{"create", template, "--"}
	args = append(args, flags...)
	args = append(args, name)
	return &task.ShellCommand{
		Bin:  p.bin,
		Args: args,
	}
}

func (p *pnpm) Install(name []string) *task.ShellCommand {
	args := []string{"install"}
	args = append(args, name...)
	return &task.ShellCommand{
		Bin:  p.bin,
		Args: args,
	}
}

func (p *pnpm) InstallDev(name []string) *task.ShellCommand {
	args := []string{"install", "-D"}
	args = append(args, name...)
	return &task.ShellCommand{
		Bin:  p.bin,
		Args: args,
	}
}

func (p *pnpm) Uninstall(name []string) *task.ShellCommand {
	args := []string{"rm"}
	args = append(args, name...)
	return &task.ShellCommand{
		Bin:  p.bin,
		Args: args,
	}
}

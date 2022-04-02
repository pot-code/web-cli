package shell

import "github.com/pot-code/web-cli/internal/task"

func GoModTidy() *task.ShellCommand {
	return &task.ShellCommand{
		Bin:  "go",
		Args: []string{"mod", "tidy"},
	}
}

func GoWire(p string) *task.ShellCommand {
	return &task.ShellCommand{
		Bin:  "wire",
		Args: []string{p},
	}
}

func GoImports(p string) *task.ShellCommand {
	return &task.ShellCommand{
		Bin:  "goimports",
		Args: []string{"-w", p},
	}
}

func GoEntInit(module string) *task.ShellCommand {
	return &task.ShellCommand{
		Bin:  "ent",
		Args: []string{"init", module},
	}
}

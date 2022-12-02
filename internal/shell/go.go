package shell

import "github.com/pot-code/web-cli/internal/task"

func GoModTidy() *task.ShellCommandTask {
	return &task.ShellCommandTask{
		Bin:  "go",
		Args: []string{"mod", "tidy"},
	}
}

func GoWire(p string) *task.ShellCommandTask {
	return &task.ShellCommandTask{
		Bin:  "wire",
		Args: []string{p},
	}
}

func GoImports(p string) *task.ShellCommandTask {
	return &task.ShellCommandTask{
		Bin:  "goimports",
		Args: []string{"-w", p},
	}
}

func GoEntInit(module string) *task.ShellCommandTask {
	return &task.ShellCommandTask{
		Bin:  "ent",
		Args: []string{"init", module},
	}
}

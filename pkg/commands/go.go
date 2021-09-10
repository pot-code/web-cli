package commands

import "github.com/pot-code/web-cli/pkg/core"

func GoModTidy() *core.ShellCommand {
	return &core.ShellCommand{
		Bin:  "go",
		Args: []string{"mod", "tidy"},
	}
}

func GoWire(p string) *core.ShellCommand {
	return &core.ShellCommand{
		Bin:  "wire",
		Args: []string{p},
	}
}

func GoImports(p string) *core.ShellCommand {
	return &core.ShellCommand{
		Bin:  "goimports",
		Args: []string{"-w", p},
	}
}

func GoEntInit(module string) *core.ShellCommand {
	return &core.ShellCommand{
		Bin:  "ent",
		Args: []string{"init", module},
	}
}

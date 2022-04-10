package pkm

import (
	"fmt"

	"github.com/pot-code/web-cli/internal/task"
)

type PackageManager interface {
	Create(template, name string, flags []string) *task.ShellCommand
	Install(name []string) *task.ShellCommand
	InstallDev(name []string) *task.ShellCommand
	Uninstall(name []string) *task.ShellCommand
}

func NewPackageManager(bin string) PackageManager {
	switch bin {
	case "pnpm":
		return newPnpm(bin)
	case "npm":
		return newNpm(bin)
	case "yarn":
		return newYarn(bin)
	}
	panic(fmt.Sprintf("unknown package manager '%s'", bin))
}

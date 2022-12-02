package pkm

import (
	"fmt"

	"github.com/pot-code/web-cli/internal/task"
)

type PackageManager interface {
	Create(template, name string, flags []string) *task.ShellCommandTask
	Install(name []string) *task.ShellCommandTask
	InstallDev(name []string) *task.ShellCommandTask
	Uninstall(name []string) *task.ShellCommandTask
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

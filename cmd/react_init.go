package cmd

import (
	"github.com/pot-code/web-cli/internal/command"
	"github.com/pot-code/web-cli/internal/pkm"
	"github.com/pot-code/web-cli/internal/shell"
	"github.com/pot-code/web-cli/internal/task"
	"github.com/pot-code/web-cli/internal/util"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type ReactInitConfig struct {
	GenType        string `flag:"type" alias:"t" usage:"project type" validate:"oneof=vanilla nextjs"`
	PackageManager string `flag:"package-manager" alias:"pm" usage:"choose package manager" validate:"oneof=pnpm npm yarn"`
	ProjectName    string `arg:"0" alias:"project_name" validate:"required,nature"`
}

var ReactInitCmd = command.NewCliCommand("init", "create react project",
	&ReactInitConfig{
		GenType:        "vanilla",
		PackageManager: "pnpm",
	},
	command.WithArgUsage("project_name"),
).AddHandlers(
	new(InitReact),
).BuildCommand()

type InitReact struct {
	ProjectName    string
	PackageManager string
}

func (ir *InitReact) Handle(c *cli.Context, cfg interface{}) error {
	config := cfg.(*ReactInitConfig)
	gt := config.GenType

	ir.ProjectName = config.ProjectName
	ir.PackageManager = config.PackageManager

	if util.Exists(ir.ProjectName) {
		log.Infof("folder '%s' already exists", ir.ProjectName)
		return nil
	}

	if gt == "vanilla" {
		return ir.vanilla().Run()
	}
	if gt == "nextjs" {
		return ir.nextjs().Run()
	}
	return nil
}

func (ir *InitReact) nextjs() task.Task {
	pm := pkm.NewPackageManager(ir.PackageManager)
	return pm.Create("next-app", ir.ProjectName, []string{"--typescript"})
}

func (ir *InitReact) vanilla() task.Task {
	return task.NewSequentialExecutor(
		[]task.Task{
			shell.GitClone("https://github.com/pot-code/react-template.git", ir.ProjectName),
			shell.GitDeleteHistory(ir.ProjectName),
		},
	)
}

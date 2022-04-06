package cmd

import (
	"bytes"

	"github.com/pot-code/web-cli/internal/command"
	"github.com/pot-code/web-cli/internal/task"
	"github.com/urfave/cli/v2"
)

type GoViperTagConfig struct {
	ConfigPath string `arg:"0" alias:"config_path" validate:"required"`
	StructName string `flag:"struct" alias:"s" usage:"struct name" validate:"required"`
}

var GoViperTagCmd = command.NewCliCommand("viper", "transform config struct to pflag",
	&GoViperTagConfig{
		StructName: "AppConfig",
	},
	command.WithAlias([]string{"v"}),
	command.WithArgUsage("config_path"),
).AddHandlers(AddViperTag).BuildCommand()

var AddViperTag = command.InlineHandler(func(c *cli.Context, cfg interface{}) error {
	config := cfg.(*GoViperTagConfig)

	var buf bytes.Buffer
	return task.NewSequentialExecutor(
		[]task.Task{
			&task.ShellCommand{
				Bin: "gomodifytags",
				Args: []string{
					"-file",
					config.ConfigPath,
					"-all",
					"-add-tags",
					"mapstructure,yaml",
				},
				Out: &buf,
			},
			&task.FileGenerator{
				Name:      config.ConfigPath,
				Overwrite: true,
				Data:      &buf,
			},
		},
	).Run()
})

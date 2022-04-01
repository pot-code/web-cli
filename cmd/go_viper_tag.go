package cmd

import (
	"bytes"

	"github.com/pot-code/web-cli/pkg/task"
	"github.com/pot-code/web-cli/pkg/util"
	"github.com/urfave/cli/v2"
)

type GoViperTagConfig struct {
	ConfigPath string `arg:"0" alias:"config_path" validate:"required"`
	StructName string `flag:"struct" alias:"s" usage:"struct name" validate:"required"`
}

var GoViperTagCmd = util.NewCliCommand("viper", "transform config struct to pflag",
	&GoViperTagConfig{
		StructName: "AppConfig",
	},
	util.WithAlias([]string{"v"}),
	util.WithArgUsage("config_path"),
).AddFeature(AddViperTag).ExportCommand()

var AddViperTag = util.NoCondFeature(func(c *cli.Context, cfg interface{}) error {
	config := cfg.(*GoViperTagConfig)

	var outData bytes.Buffer
	return util.NewTaskComposer("").AddFile(
		&task.FileDesc{
			Path:      config.ConfigPath,
			Overwrite: true,
			Source: func(buf *bytes.Buffer) error {
				buf.Write(buf.Bytes())
				return nil
			},
		},
	).AddBeforeCommand(
		&task.ShellCommand{
			Bin: "gomodifytags",
			Args: []string{
				"-file",
				config.ConfigPath,
				// "-struct",
				// config.StructName,
				"-all",
				"-add-tags",
				"mapstructure,yaml",
			},
			Out: &outData,
		},
	).Run()
})

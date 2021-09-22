package cmd

import (
	"bytes"

	"github.com/pot-code/web-cli/pkg/core"
	"github.com/pot-code/web-cli/pkg/util"
	"github.com/urfave/cli/v2"
)

type GoViperTagConfig struct {
	ConfigPath string `arg:"0" alias:"config_path" validate:"required"`
	StructName string `flag:"struct" alias:"s" usage:"struct name" validate:"required"`
}

var GoViperTagCmd = core.NewCliLeafCommand("viper", "transform config struct to pflag",
	&GoViperTagConfig{
		StructName: "AppConfig",
	},
	core.WithAlias([]string{"v"}),
	core.WithArgUsage("config_path"),
).AddService(AddViperTagService).ExportCommand()

var AddViperTagService = util.NoCondFunctionService(func(c *cli.Context, cfg interface{}) error {
	config := cfg.(*GoViperTagConfig)

	var outData bytes.Buffer
	return util.NewTaskComposer("").AddFile(
		&core.FileDesc{
			Path:      config.ConfigPath,
			Overwrite: true,
			Data: func() ([]byte, error) {
				return outData.Bytes(), nil
			},
		},
	).AddBeforeCommand(
		&core.ShellCommand{
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

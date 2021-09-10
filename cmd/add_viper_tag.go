package cmd

import (
	"bytes"

	"github.com/pot-code/web-cli/pkg/core"
	"github.com/pot-code/web-cli/pkg/util"
	"github.com/urfave/cli/v2"
)

type AddViperTagConfig struct {
	ConfigPath string `arg:"0" alias:"config_path" validate:"required"`
	StructName string `flag:"struct" alias:"s" usage:"struct name" validate:"required"`
}

var AddViperTagCmd = core.NewCliLeafCommand("viper", "transform config struct to pflag",
	&AddViperTagConfig{
		StructName: "AppConfig",
	},
	core.WithAlias([]string{"V"}),
	core.WithArgUsage("config_path"),
).AddService(AddViperTagService).ExportCommand()

var AddViperTagService = util.NoCondFunctionService(func(c *cli.Context, cfg interface{}) error {
	config := cfg.(*AddViperTagConfig)

	var outData bytes.Buffer
	return util.NewTaskComposer("").AddFile(
		&core.FileDesc{
			Path:      config.ConfigPath,
			Overwrite: true,
			Data: func() []byte {
				return outData.Bytes()
			},
		},
	).AddCommand(
		&core.ShellCommand{
			Bin:    "gomodifytags",
			Before: true,
			Args: []string{
				"-file",
				config.ConfigPath,
				"-struct",
				config.StructName,
				"-add-tags",
				"mapstructure,yaml",
			},
			Out: &outData,
		},
	).Run()
})

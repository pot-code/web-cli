package cmd

import (
	"bytes"

	"github.com/pot-code/web-cli/pkg/core"
	"github.com/pot-code/web-cli/pkg/util"
	"github.com/urfave/cli/v2"
)

type addViperTagConfig struct {
	ConfigPath string `arg:"0" alias:"config_path" validate:"required"`
	StructName string `flag:"struct" validate:"required"`
}

var addViperTagCmd = &cli.Command{
	Name:      "viper",
	Aliases:   []string{"v"},
	Usage:     "transform config struct to pflags",
	ArgsUsage: "config_path",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "struct",
			Aliases: []string{"s"},
			Usage:   "struct name",
			Value:   "AppConfig",
		},
	},
	Action: func(c *cli.Context) error {
		config := new(addViperTagConfig)
		err := util.ParseConfig(c, config)
		if err != nil {
			if _, ok := err.(*util.CommandError); ok {
				cli.ShowCommandHelp(c, c.Command.Name)
			}
			return err
		}

		cmd, err := addViperFlag(config)
		if err != nil {
			return err
		}
		return cmd.Run()
	},
}

func addViperFlag(config *addViperTagConfig) (core.Runner, error) {
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
		&core.Command{
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
	), nil
}

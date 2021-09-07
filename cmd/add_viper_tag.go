package cmd

import (
	"bytes"

	"github.com/pot-code/web-cli/pkg/core"
	"github.com/pot-code/web-cli/pkg/util"
	"github.com/urfave/cli/v2"
)

type addViperTagConfig struct {
	ConfigPath string `arg:"0" name:"CONFIG_PATH" validate:"required"`
	OutName    string `name:"out"`
}

var addViperTagCmd = &cli.Command{
	Name:      "viper",
	Aliases:   []string{"v"},
	Usage:     "generate flags registration go file",
	ArgsUsage: "CONFIG_PATH",
	Flags:     []cli.Flag{
		// &cli.StringFlag{
		// 	Name:    "out",
		// 	Aliases: []string{"o"},
		// 	Usage:   "output file name",
		// 	Value:   "config_gen.go",
		// },
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
		}).AddCommand(&core.Command{
		Bin:    "gomodifytags",
		Before: true,
		Args: []string{
			"-file",
			config.ConfigPath,
			"-struct",
			"AppConfig",
			"-add-tags",
			"mapstructure,yaml",
		},
		Out: &outData,
	}), nil
}

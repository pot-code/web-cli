package cmd

import (
	"bytes"
	"fmt"
	"go/format"
	"path"

	"github.com/pot-code/web-cli/pkg/constants"
	"github.com/pot-code/web-cli/pkg/core"
	"github.com/pot-code/web-cli/pkg/util"
	"github.com/pot-code/web-cli/templates"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type genFlagsConfig struct {
	ConfigPath string `arg:"0" alias:"CONFIG_PATH" validate:"required"`
	FileName   string `flag:"name" validate:"var"`
}

var genFlagsCmd = &cli.Command{
	Name:      "flags",
	Aliases:   []string{"f"},
	Usage:     "generate flags registration go file",
	ArgsUsage: "CONFIG_PATH",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "name",
			Usage:   "generated file name",
			Aliases: []string{"n"},
			Value:   "flags_gen",
		},
	},
	Action: func(c *cli.Context) error {
		config := new(genFlagsConfig)
		err := util.ParseConfig(c, config)
		if err != nil {
			if _, ok := err.(*util.CommandError); ok {
				cli.ShowCommandHelp(c, c.Command.Name)
			}
			return err
		}

		cmd, err := generateFlagsFile(config)
		if err != nil {
			return err
		}
		return cmd.Run()
	},
}

func generateFlagsFile(config *genFlagsConfig) (core.Runner, error) {
	visitor, err := parseConfigFile(config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pflags: %w", err)
	}

	fm := []map[string]interface{}{}
	for _, fd := range visitor.fds {
		fm = append(fm, map[string]interface{}{
			"FlagType":     fd.FlagType,
			"Name":         fd.TagMeta.MarshalName,
			"DefaultValue": fd.TagMeta.Default,
			"Required":     fd.TagMeta.Required,
			"Help":         fd.Help,
		})
	}

	pkg := visitor.pkg
	fileName := config.FileName + constants.GoSuffix
	return util.NewTaskComposer(pkg).AddFile(
		&core.FileDesc{
			Path:      fileName,
			Overwrite: true,
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoGenPflags(&buf, pkg, fm)
				data, err := format.Source(buf.Bytes())
				if err != nil {
					log.WithField("err", err).Error("failed to format code")
					return buf.Bytes()
				}
				return data
			},
		},
	).AddCommand(
		&core.Command{
			Bin:  "goimports",
			Args: []string{"-w", path.Join(pkg, fileName)},
		},
		&core.Command{
			Bin:  "go",
			Args: []string{"mod", "tidy"},
		},
	), nil
}

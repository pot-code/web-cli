package cmd

import (
	"bytes"
	"fmt"
	"path"

	"github.com/pkg/errors"
	"github.com/pot-code/web-cli/internal/command"
	"github.com/pot-code/web-cli/internal/constant"
	"github.com/pot-code/web-cli/internal/shell"
	"github.com/pot-code/web-cli/internal/task"
	"github.com/pot-code/web-cli/internal/transformer"
	"github.com/pot-code/web-cli/templates"
	"github.com/urfave/cli/v2"
)

type GoFlagsConfig struct {
	ConfigPath  string `arg:"0" alias:"CONFIG_PATH" validate:"required"`
	OutFileName string `flag:"name" alias:"n" usage:"generated file name" validate:"required,var"`
}

var GoFlagsCmd = command.NewCliCommand("flags", "generate pflags registration based on struct",
	&GoFlagsConfig{
		OutFileName: "config_gen",
	},
	command.WithAlias([]string{"f"}),
	command.WithArgUsage("CONFIG_PATH"),
).AddFeature(GenViperFlags).ExportCommand()

var GenViperFlags = command.NoCondFeature(func(c *cli.Context, cfg interface{}) error {
	config := cfg.(*GoFlagsConfig)
	visitor, err := parseConfigFile(config)
	if err != nil {
		return errors.WithStack(fmt.Errorf("failed to parse pflags: %w", err))
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
	out := fmt.Sprintf("%s.%s", config.OutFileName, constant.GoSuffix)
	return task.NewSequentialExecutor(
		[]task.Task{
			(&task.FileGenerator{
				Name:      out,
				Overwrite: true,
				Data:      bytes.NewBufferString(templates.GoGenPflags(pkg, fm)),
			}).UseTransformers(transformer.GoFormatSource()),
			shell.GoImports(path.Join(pkg, out)),
			shell.GoModTidy(),
		},
	).Run()
})

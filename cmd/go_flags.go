package cmd

import (
	"bytes"
	"fmt"
	"path"

	"github.com/pkg/errors"
	"github.com/pot-code/web-cli/internal/commands"
	"github.com/pot-code/web-cli/internal/constants"
	"github.com/pot-code/web-cli/internal/task"
	"github.com/pot-code/web-cli/internal/transformation"
	"github.com/pot-code/web-cli/internal/util"
	"github.com/pot-code/web-cli/templates"
	"github.com/urfave/cli/v2"
)

type GoFlagsConfig struct {
	ConfigPath  string `arg:"0" alias:"CONFIG_PATH" validate:"required"`
	OutFileName string `flag:"name" alias:"n" usage:"generated file name" validate:"required,var"`
}

var GoFlagsCmd = util.NewCliCommand("flags", "generate pflags registration based on struct",
	&GoFlagsConfig{
		OutFileName: "config_gen",
	},
	util.WithAlias([]string{"f"}),
	util.WithArgUsage("CONFIG_PATH"),
).AddFeature(GenViperFlags).ExportCommand()

var GenViperFlags = util.NoCondFeature(func(c *cli.Context, cfg interface{}) error {
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
	out := fmt.Sprintf("%s.%s", config.OutFileName, constants.GoSuffix)
	return util.NewTaskComposer(pkg).AddFile(
		&task.FileDesc{
			Path:      out,
			Overwrite: true,
			Source: func(buf *bytes.Buffer) error {
				templates.WriteGoGenPflags(buf, pkg, fm)
				return nil
			},
			Transforms: []task.Pipeline{transformation.GoFormatSource},
		},
	).AddCommand(
		commands.GoImports(path.Join(pkg, out)),
		commands.GoModTidy(),
	).Run()
})

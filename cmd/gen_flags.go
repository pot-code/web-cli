package cmd

import (
	"bytes"
	"fmt"
	"path"

	"github.com/pkg/errors"
	"github.com/pot-code/web-cli/pkg/constants"
	"github.com/pot-code/web-cli/pkg/core"
	"github.com/pot-code/web-cli/pkg/util"
	"github.com/pot-code/web-cli/templates"
	"github.com/urfave/cli/v2"
)

type GenFlagsConfig struct {
	ConfigPath string `arg:"0" alias:"CONFIG_PATH" validate:"required"`
	FileName   string `flag:"name" alias:"n" usage:"generated file name" validate:"required,var"`
}

var GenFlagsCmd = core.NewCliLeafCommand("flags", "generate flags registration go file",
	&GenFlagsConfig{
		FileName: "config_gen",
	},
	core.WithAlias([]string{"f"}),
	core.WithArgUsage("CONFIG_PATH"),
).AddService(GenViperFlagsService).ExportCommand()

var GenViperFlagsService = util.NoCondFunctionService(func(c *cli.Context, cfg interface{}) error {
	config := cfg.(*GenFlagsConfig)
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
	fileName := config.FileName + constants.GoSuffix
	return util.NewTaskComposer(pkg).AddFile(
		&core.FileDesc{
			Path:      fileName,
			Overwrite: true,
			Data: func() ([]byte, error) {
				var buf bytes.Buffer

				templates.WriteGoGenPflags(&buf, pkg, fm)
				// TODO: format output
				return buf.Bytes(), nil
			},
		},
	).AddCommand(
		&core.ShellCommand{
			Bin:  "goimports",
			Args: []string{"-w", path.Join(pkg, fileName)},
		},
		&core.ShellCommand{
			Bin:  "go",
			Args: []string{"mod", "tidy"},
		},
	).Run()
})

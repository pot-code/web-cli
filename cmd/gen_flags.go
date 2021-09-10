package cmd

import (
	"bytes"
	"fmt"
	"go/format"
	"path"

	"github.com/pkg/errors"
	"github.com/pot-code/web-cli/pkg/constants"
	"github.com/pot-code/web-cli/pkg/core"
	"github.com/pot-code/web-cli/pkg/util"
	"github.com/pot-code/web-cli/templates"
	log "github.com/sirupsen/logrus"
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
).AddService(new(GenViperFlagsService)).ExportCommand()

type GenViperFlagsService struct{}

var _ core.CommandService = &GenGolangBeService{}

func (gvf *GenViperFlagsService) Cond(c *cli.Context) bool {
	return true
}

func (gvf *GenViperFlagsService) Handle(c *cli.Context, cfg interface{}) error {
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
	).Run()
}

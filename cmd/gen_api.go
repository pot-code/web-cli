package cmd

import (
	"bytes"
	"os"
	"path"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/pot-code/web-cli/core"
	"github.com/pot-code/web-cli/templates"
	"github.com/pot-code/web-cli/util"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const cmdApiName = "api"

type genApiConfig struct {
	GenType     string `name:"type" validate:"required,oneof=go"`    // generation type
	PackageName string `arg:"0" name:"NAME" validate:"required,var"` // go pkg name
	ProjectName string
	Author      string
	Model       string
	Root        string `name:"root" validate:"required,dir"` // path root under which to generate api
}

var genAPICmd = &cli.Command{
	Name:      cmdApiName,
	Usage:     "generate an api producer",
	ArgsUsage: "NAME",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "type",
			Aliases: []string{"T"},
			Usage:   "api type (go)",
			Value:   "go",
		},
		&cli.StringFlag{
			Name:    "root",
			Aliases: []string{"r"},
			Usage:   "generation root",
			Value:   "api",
		},
	},
	Action: func(c *cli.Context) error {
		config := new(genApiConfig)
		err := util.ParseConfig(c, config)
		if err != nil {
			if _, ok := err.(*util.CommandError); ok {
				cli.ShowCommandHelp(c, cmdApiName)
			}
			return err
		}

		name := strings.ToLower(config.PackageName)
		name = strings.ReplaceAll(name, "-", "_")
		log.Debug("transformed module name: ", name)
		config.PackageName = name

		if config.GenType == "go" {
			meta, err := util.ParseGoMod("go.mod")
			if err != nil {
				return err
			}

			config.Model = strcase.ToCamel(name)
			config.Author = meta.Author
			config.ProjectName = meta.ProjectName

			dir := path.Join(config.Root, config.PackageName)
			_, err = os.Stat(dir)
			if err == nil {
				log.Infof("[skipped]'%s' already exists", dir)
				return nil
			}

			gen := newGolangApiGenerator(config)
			if err := gen.Run(); err != nil {
				gen.Cleanup()
				return err
			}
		}
		return nil
	},
}

func newGolangApiGenerator(config *genApiConfig) core.Generator {
	return util.NewTaskComposer(
		path.Join(config.Root, config.PackageName),
		&core.FileDesc{
			Path: "transport/http.go",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoApiHttp(&buf, config.ProjectName, config.Author, config.PackageName, config.Model)
				return buf.Bytes()
			},
		},
		&core.FileDesc{
			Path: "model.go",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoApiModel(&buf, config.PackageName, config.Model)
				return buf.Bytes()
			},
		},
		&core.FileDesc{
			Path: "repo.go",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoApiRepo(&buf, config.PackageName, config.Model)
				return buf.Bytes()
			},
		},
		&core.FileDesc{
			Path: "service.go",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoApiService(&buf, config.PackageName, config.Model)
				return buf.Bytes()
			},
		},
	)
}

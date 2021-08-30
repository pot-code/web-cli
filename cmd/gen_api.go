package cmd

import (
	"bytes"
	"os"
	"path"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/pot-code/web-cli/pkg/core"
	"github.com/pot-code/web-cli/pkg/util"
	"github.com/pot-code/web-cli/templates"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type genApiConfig struct {
	GenType     string `name:"type" validate:"required,oneof=go"`    // generation type
	ModuleName  string `arg:"0" name:"NAME" validate:"required,var"` // go pkg name
	PackagePath string
	ProjectName string
	ModelName   string
	Author      string
	Root        string `name:"root" validate:"required"` // path root under which to generate api
}

var genAPICmd = &cli.Command{
	Name:      "api",
	Usage:     "generate an api module",
	ArgsUsage: "MODULE_NAME",
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
			Usage:   "root directory",
			Value:   "api",
		},
	},
	Action: func(c *cli.Context) error {
		config := new(genApiConfig)
		err := util.ParseConfig(c, config)
		if err != nil {
			if _, ok := err.(*util.CommandError); ok {
				cli.ShowCommandHelp(c, c.Command.Name)
			}
			return err
		}

		pkgName := strings.ReplaceAll(config.ModuleName, "-", "_")
		pkgName = strcase.ToSnake(pkgName)
		config.PackagePath = pkgName

		if config.GenType == "go" {
			meta, err := util.ParseGoMod("go.mod")
			if err != nil {
				return err
			}

			config.ModelName = strcase.ToCamel(pkgName)
			log.Debug("generated module name: ", config.ModelName)
			config.Author = meta.Author
			config.ProjectName = meta.ProjectName

			dir := path.Join(config.Root, config.PackagePath)
			log.Debug("output path: ", dir)
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
		path.Join(config.Root, config.PackagePath),
		&core.FileDesc{
			Path: "transport/http.go",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoApiHttp(&buf, config.ProjectName, config.Author, config.PackagePath, config.ModelName)
				return buf.Bytes()
			},
		},
		&core.FileDesc{
			Path: "model.go",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoApiModel(&buf, config.PackagePath, config.ModelName)
				return buf.Bytes()
			},
		},
		&core.FileDesc{
			Path: "repo.go",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoApiRepo(&buf, config.PackagePath, config.ModelName)
				return buf.Bytes()
			},
		},
		&core.FileDesc{
			Path: "service.go",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoApiService(&buf, config.PackagePath, config.ModelName)
				return buf.Bytes()
			},
		},
	)
}

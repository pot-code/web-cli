package cmd

import (
	"bytes"
	"fmt"
	"os"
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

			dir := config.PackagePath
			log.Debug("output path: ", dir)
			_, err = os.Stat(dir)
			if err == nil {
				log.Infof("[skipped]'%s' already exists", dir)
				return nil
			}

			gen := generateGoApi(config)
			if err := gen.Run(); err != nil {
				return err
			}
		}
		return nil
	},
}

func generateGoApi(config *genApiConfig) core.Runner {
	return util.NewTaskComposer("",
		&core.FileDesc{
			Path: fmt.Sprintf("%s/%s_handler.go", "server", config.PackagePath),
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoApiHandler(&buf, config.ProjectName, config.Author, config.PackagePath, config.ModelName)
				return buf.Bytes()
			},
		},
		&core.FileDesc{
			Path: fmt.Sprintf("%s/model.go", config.PackagePath),
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoApiModel(&buf, config.PackagePath, config.ModelName)
				return buf.Bytes()
			},
		},
		// &core.FileDesc{
		// 	Path: fmt.Sprintf("%s/repo.go", config.PackagePath),
		// 	Data: func() []byte {
		// 		var buf bytes.Buffer

		// 		templates.WriteGoApiRepo(&buf, config.PackagePath, config.ModelName)
		// 		return buf.Bytes()
		// 	},
		// },
		&core.FileDesc{
			Path: fmt.Sprintf("%s/service.go", config.PackagePath),
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoApiService(&buf, config.PackagePath, config.ModelName)
				return buf.Bytes()
			},
		},
	).AddCommand(&core.Command{
		Bin:    "ent",
		Args:   []string{"init", config.ModelName},
		Before: true,
	})
}

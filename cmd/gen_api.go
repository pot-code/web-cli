package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/pot-code/web-cli/core"
	"github.com/pot-code/web-cli/templates"
	"github.com/pot-code/web-cli/util"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const cmdApiName = "api"

type genApiConfig struct {
	GenType     string // generation type
	PackageName string // go pkg name
	Project     string
	Author      string
	Model       string
	Root        string // path root under which to generate api
}

var genAPICmd = &cli.Command{
	Name:      cmdApiName,
	Usage:     "generate an api producer",
	ArgsUsage: "NAME",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "type",
			Aliases: []string{"T"},
			Usage:   "api type (nest/go)",
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
		config, err := getGenApiConfig(c)
		if err != nil {
			if _, ok := err.(*util.CommandError); ok {
				cli.ShowCommandHelp(c, cmdApiName)
			}
			return err
		}

		if config.GenType == "go" {
			dir := path.Join(config.Root, config.PackageName)
			_, err := os.Stat(dir)
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

				templates.WriteGoApiHttp(&buf, config.Project, config.Author, config.PackageName, config.Model)
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

func getGenApiConfig(c *cli.Context) (*genApiConfig, error) {
	config := &genApiConfig{}
	name := c.Args().Get(0) // module name
	if name == "" {
		return nil, util.NewCommandError(cmdBackendName, fmt.Errorf("NAME must be specified"))
	}

	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, "-", "_")
	log.Debug("transformed module name: ", name)
	if err := util.ValidateVarName(name); err != nil {
		return nil, util.NewCommandError(cmdBackendName, errors.Wrap(err, "invalid NAME"))
	}

	genType := strings.ToLower(c.String("type"))
	if genType == "" {
		return nil, util.NewCommandError(cmdBackendName, fmt.Errorf("type is empty"))
	} else if genType == "go" {
		log.Debugf("opening '%s'", "go.mod")
		meta, err := util.ParseGoMod("go.mod")
		if err != nil {
			return nil, err
		}

		config.Model = strcase.ToCamel(name)
		config.Author = meta.Author
		config.Project = meta.ProjectName
	} else if genType == "nest" {
		panic("not implemented") // TODO:
	} else {
		return nil, util.NewCommandError(cmdBackendName, fmt.Errorf("unsupported type %s", genType))
	}

	root := c.String("root")

	config.PackageName = name
	config.GenType = genType
	config.Root = root
	log.WithFields(log.Fields{"author": config.Author, "project": config.Project, "model": config.Model, "type": genType, "root": root}).Debugf("parsed meta data")
	return config, nil
}

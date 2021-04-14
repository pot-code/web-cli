package cmd

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/pot-code/web-cli/generate"
	"github.com/pot-code/web-cli/util"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const CmdApiName = "api"

type GenApiConfig struct {
	GenType     string // generation type
	PackageName string // go pkg name
	Project     string
	Author      string
	Model       string
	Root        string // path root under which to generate api
}

var genAPICmd = &cli.Command{
	Name:      CmdApiName,
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
				cli.ShowCommandHelp(c, CmdApiName)
			}
			return err
		}

		if config.GenType == "go" {
			gen := generate.NewGolangApiGenerator(&generate.GolangApiConfig{
				PackageName: config.PackageName,
				Project:     config.Project,
				Author:      config.Author,
				Model:       config.Model,
				Root:        config.Root,
			})
			if err := gen.Gen(); err != nil {
				gen.Cleanup()
				return err
			}
		}
		return nil
	},
}

func getGenApiConfig(c *cli.Context) (*GenApiConfig, error) {
	config := &GenApiConfig{}
	name := c.Args().Get(0) // module name
	if name == "" {
		return nil, util.NewCommandError(CmdBackendName, fmt.Errorf("NAME must be specified"))
	}

	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, "-", "_")
	log.Debug("transformed module name: ", name)
	if err := util.ValidateProjectName(name); err != nil {
		return nil, util.NewCommandError(CmdBackendName, errors.Wrap(err, "invalid NAME"))
	}

	genType := strings.ToLower(c.String("type"))
	if genType == "" {
		return nil, util.NewCommandError(CmdBackendName, fmt.Errorf("type is empty"))
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
		return nil, util.NewCommandError(CmdBackendName, fmt.Errorf("unsupported type %s", genType))
	}

	root := c.String("root")

	config.PackageName = name
	config.GenType = genType
	config.Root = root
	log.WithFields(log.Fields{"author": config.Author, "project": config.Project, "model": config.Model, "type": genType, "root": root}).Debugf("parsed meta data")
	return config, nil
}

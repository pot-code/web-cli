package cmd

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/pot-code/web-cli/generate"
	"github.com/pot-code/web-cli/util"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const cmdBackendName = "backend"

type genBEConfig struct {
	GenType     string // generation type
	ProjectName string // project name
	Author      string // project author name
	Version     string // version number
}

var generateBECmd = &cli.Command{
	Name:      cmdBackendName,
	Aliases:   []string{"be"},
	Usage:     "generate backends",
	ArgsUsage: "NAME",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "type",
			Aliases: []string{"t"},
			Usage:   "backend type (nest/go)",
			Value:   "go",
		},
		&cli.StringFlag{
			Name:     "author",
			Aliases:  []string{"a"},
			Usage:    "author name for the app (required)",
			Required: true,
		},
		&cli.StringFlag{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "version number for go.mod generation",
			Value:   "1.14",
		},
	},
	Action: func(c *cli.Context) error {
		config, err := getGenBEConfig(c)
		if err != nil {
			if _, ok := err.(*util.CommandError); ok {
				cli.ShowCommandHelp(c, cmdBackendName)
			}
			return err
		}

		if config.GenType == "go" {
			gen := generate.NewGolangBackendGenerator(&generate.GolangBackendConfig{
				ProjectName: config.ProjectName,
				Author:      config.Author,
				Version:     config.Version,
			})
			if err := gen.Run(); err != nil {
				gen.Cleanup()
				return err
			}
		}
		return nil
	},
}

func getGenBEConfig(c *cli.Context) (*genBEConfig, error) {
	name := c.Args().Get(0)
	if name == "" {
		return nil, util.NewCommandError(cmdBackendName, fmt.Errorf("NAME must be specified"))
	}
	if err := util.ValidateVarName(name); err != nil {
		return nil, util.NewCommandError(cmdBackendName, errors.Wrap(err, "invalid NAME"))
	}

	author := c.String("author")
	version := c.String("version")
	genType := strings.ToLower(c.String("type"))
	if genType == "" {
		return nil, util.NewCommandError(cmdBackendName, fmt.Errorf("type is empty"))
	} else if genType == "go" {
		if author == "" {
			return nil, util.NewCommandError(cmdBackendName, fmt.Errorf("author is empty"))
		}
		if err := util.ValidateUserName(author); err != nil {
			return nil, util.NewCommandError(cmdBackendName, errors.Wrap(err, "invalid author name"))
		}
		if version == "" {
			return nil, util.NewCommandError(cmdBackendName, fmt.Errorf("version is empty when type is 'go'"))
		}
		if err := util.ValidateVersion(version); err != nil {
			return nil, util.NewCommandError(cmdBackendName, errors.Wrap(err, "invalid version"))
		}
	} else if genType == "nest" {
		panic("not implemented") // TODO:
	} else {
		return nil, util.NewCommandError(cmdBackendName, fmt.Errorf("unsupported type %s", genType))
	}

	log.WithFields(log.Fields{"author": author, "project": name, "version": version, "type": genType}).Debugf("parsed meta data")
	return &genBEConfig{
		GenType:     genType,
		ProjectName: name,
		Author:      author,
		Version:     version,
	}, nil
}
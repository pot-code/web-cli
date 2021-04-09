package cmd

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/pot-code/web-cli/util"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const CmdBackendName = "backend"

type GenBEConfig struct {
	GenType     string // generation type
	ProjectName string // project name
	Author      string // project author name
	Version     string // version number
}

var genBE = &cli.Command{
	Name:      CmdBackendName,
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
				cli.ShowCommandHelp(c, CmdBackendName)
			}
			return err
		}

		if config.GenType == "go" {
			if err := NewGolangBackendGenerator(config).Gen(); err != nil {
				return err
			}
		}
		return nil
	},
}

func getGenBEConfig(c *cli.Context) (*GenBEConfig, error) {
	name := c.Args().Get(0)
	if name == "" {
		return nil, util.NewCommandError(CmdBackendName, fmt.Errorf("NAME must be specified"))
	}
	if err := util.ValidateProjectName(name); err != nil {
		return nil, util.NewCommandError(CmdBackendName, errors.Wrap(err, "invalid NAME"))
	}

	genType := strings.ToLower(c.String("type"))
	author := c.String("author")
	version := c.String("version")

	if genType == "" {
		return nil, util.NewCommandError(CmdBackendName, fmt.Errorf("type is empty"))
	} else if genType == "go" {
		if author == "" {
			return nil, util.NewCommandError(CmdBackendName, fmt.Errorf("author is empty"))
		}
		if err := util.ValidateUserName(author); err != nil {
			return nil, util.NewCommandError(CmdBackendName, errors.Wrap(err, "invalid author name"))
		}
	} else if genType == "nest" {
		panic("not implemented") // TODO:
	} else {
		return nil, util.NewCommandError(CmdBackendName, fmt.Errorf("unsupported type %s", genType))
	}

	if genType == "go" {
		if version == "" {
			return nil, util.NewCommandError(CmdBackendName, fmt.Errorf("version is empty when type is 'go'"))
		}
		if err := util.ValidateVersion(version); err != nil {
			return nil, util.NewCommandError(CmdBackendName, errors.Wrap(err, "invalid version"))
		}
	}

	log.WithFields(log.Fields{"author": author, "project": name, "version": version, "type": genType}).Debugf("parsed meta data")
	return &GenBEConfig{
		GenType:     genType,
		ProjectName: name,
		Author:      author,
		Version:     version,
	}, nil
}

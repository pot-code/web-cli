package cmd

import (
	"bytes"
	"os"

	"github.com/pot-code/web-cli/pkg/core"
	"github.com/pot-code/web-cli/pkg/util"
	"github.com/pot-code/web-cli/templates"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const cmdBackendName = "backend"

type genBEConfig struct {
	GenType     string `name:"type" validate:"required,oneof=go"`    // generation type
	ProjectName string `arg:"0" name:"NAME" validate:"required,var"` // project name
	Author      string `name:"author" validate:"required,var"`       // project author name
	Version     string `name:"version" validate:"required,version"`  // version number
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
			Usage:   "backend type (go)",
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
		config := new(genBEConfig)
		err := util.ParseConfig(c, config)
		if err != nil {
			if _, ok := err.(*util.CommandError); ok {
				cli.ShowCommandHelp(c, cmdBackendName)
			}
			return err
		}

		if config.GenType == "go" {
			_, err := os.Stat(config.ProjectName)
			if err == nil {
				log.Infof("[skipped]'%s' already exists", config.ProjectName)
				return nil
			}

			gen := newGolangBackendGenerator(config)
			if err := gen.Run(); err != nil {
				gen.Cleanup()
				return err
			}
		}
		return nil
	},
}

func newGolangBackendGenerator(config *genBEConfig) core.Generator {
	return util.NewTaskComposer(
		config.ProjectName,
		&core.FileDesc{
			Path: "cmd/http.go",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoBackendCmdHttp(&buf, config.ProjectName, config.Author)
				return buf.Bytes()
			},
		},
		&core.FileDesc{
			Path: "config/def.go",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoBackendConfig(&buf, config.ProjectName)
				return buf.Bytes()
			},
		},
		&core.FileDesc{
			Path: "main.go",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoBackendMain(&buf, config.ProjectName, config.Author)
				return buf.Bytes()
			},
		},
		&core.FileDesc{
			Path: "go.mod",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoBackendMod(&buf, config.ProjectName, config.Author, config.Version)
				return buf.Bytes()
			},
		},
	)
}

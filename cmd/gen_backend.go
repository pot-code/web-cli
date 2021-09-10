package cmd

import (
	"bytes"
	"os"
	"path"

	"github.com/pot-code/web-cli/pkg/core"
	"github.com/pot-code/web-cli/pkg/util"
	"github.com/pot-code/web-cli/templates"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type genBEConfig struct {
	GenType     string `flag:"type" validate:"required,oneof=go"`             // generation type
	ProjectName string `arg:"0" alias:"project_name" validate:"required,var"` // project name
	Author      string `flag:"author" validate:"required,var"`                // project author name
	Version     string `flag:"version" validate:"required,version"`           // version number
}

var generateBECmd = &cli.Command{
	Name:      "backend",
	Aliases:   []string{"be"},
	Usage:     "generate backends",
	ArgsUsage: "project_name",
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
			Value:   "1.16",
		},
	},
	Action: func(c *cli.Context) error {
		config := new(genBEConfig)
		err := util.ParseConfig(c, config)
		if err != nil {
			if _, ok := err.(*util.CommandError); ok {
				cli.ShowCommandHelp(c, c.Command.Name)
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
				return err
			}
		}
		return nil
	},
}

func newGolangBackendGenerator(config *genBEConfig) core.Runner {
	return util.NewTaskComposer(config.ProjectName).AddFile(
		&core.FileDesc{
			Path: "cmd/web.go",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoBackendCmdWeb(&buf, config.ProjectName, config.Author)
				return buf.Bytes()
			},
		},
		&core.FileDesc{
			Path: "bootstrap/config.go",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoBackendBootstrapConfig(&buf, config.ProjectName)
				return buf.Bytes()
			},
		},
		&core.FileDesc{
			Path: "bootstrap/create.go",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoBackendBootstrapCreate(&buf, config.ProjectName, config.Author)
				return buf.Bytes()
			},
		},
		&core.FileDesc{
			Path: "server/routes.go",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoBackendServerRoutes(&buf, config.ProjectName, config.Author)
				return buf.Bytes()
			},
		},
		&core.FileDesc{
			Path: "server/server.go",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoBackendServerServer(&buf)
				return buf.Bytes()
			},
		},
		&core.FileDesc{
			Path: "server/wire.go",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoBackendServerWire(&buf, config.ProjectName, config.Author)
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
		&core.FileDesc{
			Path: ".vscode/settings.json",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoBackendVscodeSettings(&buf)
				return buf.Bytes()
			},
		},
		&core.FileDesc{
			Path: "Dockerfile",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoBackendDockerfile(&buf)
				return buf.Bytes()
			},
		},
		&core.FileDesc{
			Path: "Makefile",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoMakefile(&buf)
				return buf.Bytes()
			},
		},
		&core.FileDesc{
			Path: "air.toml",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoAirConfig(&buf)
				return buf.Bytes()
			},
		},
	).AddCommand(
		&core.Command{
			Bin:  "go",
			Args: []string{"mod", "tidy"},
			Dir:  path.Join("./" + config.ProjectName),
		},
		&core.Command{
			Bin:  "wire",
			Args: []string{"./server"},
			Dir:  path.Join("./" + config.ProjectName),
		},
	)
}

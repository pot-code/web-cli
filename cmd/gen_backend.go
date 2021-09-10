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

type GenBeConfig struct {
	GenType     string `flag:"type" alias:"t" usage:"backend type" validate:"required,oneof=go"`
	ProjectName string `arg:"0" alias:"project_name" validate:"required,var"`
	AuthorName  string `flag:"author" alias:"a" usage:"author name for the app" validate:"required,var"`
	GoVersion   string `flag:"version" alias:"v" usage:"specify go version" validate:"required,version"`
}

var GenBeCmd = core.NewCliLeafCommand("backend", "generate backends",
	&GenBeConfig{
		GenType:   "go",
		GoVersion: "1.16",
	},
	core.WithArgUsage("project_name"),
	core.WithAlias([]string{"be"}),
).AddService(new(GenGolangBeService)).ExportCommand()

type GenGolangBeService struct{}

var _ core.CommandService = &GenGolangBeService{}

func (ggb *GenGolangBeService) Cond(c *cli.Context) bool {
	return c.String("type") == "go"
}

func (ggb *GenGolangBeService) Handle(c *cli.Context, cfg interface{}) error {
	config := cfg.(*GenBeConfig)
	projectName := config.ProjectName
	authorName := config.AuthorName

	_, err := os.Stat(projectName)
	if err == nil {
		log.Infof("[skipped]'%s' already exists", projectName)
		return nil
	}

	return util.NewTaskComposer(projectName).AddFile(
		&core.FileDesc{
			Path: "cmd/web.go",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoBackendCmdWeb(&buf, projectName, authorName)
				return buf.Bytes()
			},
		},
		&core.FileDesc{
			Path: "bootstrap/config.go",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoBackendBootstrapConfig(&buf, projectName)
				return buf.Bytes()
			},
		},
		&core.FileDesc{
			Path: "bootstrap/create.go",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoBackendBootstrapCreate(&buf, projectName, authorName)
				return buf.Bytes()
			},
		},
		&core.FileDesc{
			Path: "server/routes.go",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoBackendServerRoutes(&buf, projectName, authorName)
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

				templates.WriteGoBackendServerWire(&buf, projectName, authorName)
				return buf.Bytes()
			},
		},
		&core.FileDesc{
			Path: "main.go",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoBackendMain(&buf, projectName, authorName)
				return buf.Bytes()
			},
		},
		&core.FileDesc{
			Path: "go.mod",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoBackendMod(&buf, projectName, authorName, config.GoVersion)
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
			Dir:  path.Join("./" + projectName),
		},
		&core.Command{
			Bin:  "wire",
			Args: []string{"./server"},
			Dir:  path.Join("./" + projectName),
		},
	).Run()
}

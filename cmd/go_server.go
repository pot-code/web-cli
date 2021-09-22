package cmd

import (
	"bytes"
	"os"
	"path"

	"github.com/pot-code/web-cli/pkg/core"
	"github.com/pot-code/web-cli/pkg/transform"
	"github.com/pot-code/web-cli/pkg/util"
	"github.com/pot-code/web-cli/templates"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type GoServerConfig struct {
	GenType     string `flag:"type" alias:"t" usage:"backend type" validate:"required,oneof=go"`
	ProjectName string `arg:"0" alias:"project_name" validate:"required,var"`
	AuthorName  string `flag:"author" alias:"a" usage:"author name for the app" validate:"required,var"`
	GoVersion   string `flag:"version" alias:"v" usage:"specify go version" validate:"required,version"`
}

var GoServerCmd = core.NewCliLeafCommand("server", "generate backends",
	&GoServerConfig{
		GenType:   "go",
		GoVersion: "1.16",
	},
	core.WithArgUsage("project_name"),
	core.WithAlias([]string{"s"}),
).AddService(new(GenGolangBeService)).ExportCommand()

type GenGolangBeService struct{}

var _ core.CommandService = &GenGolangBeService{}

func (ggb *GenGolangBeService) Cond(c *cli.Context) bool {
	return c.String("type") == "go"
}

func (ggb *GenGolangBeService) Handle(c *cli.Context, cfg interface{}) error {
	config := cfg.(*GoServerConfig)
	projectName := config.ProjectName
	authorName := config.AuthorName

	_, err := os.Stat(projectName)
	if err == nil {
		log.Infof("[skipped]'%s' already exists", projectName)
		return nil
	}

	return util.NewTaskComposer(projectName).AddFile(
		&core.FileDesc{
			Path: path.Join("cmd", "web.go"),
			Data: func() ([]byte, error) {
				var buf bytes.Buffer

				templates.WriteGoBackendCmdWeb(&buf, projectName, authorName)
				return buf.Bytes(), nil
			},
			Transforms: []core.Transform{transform.GoFormatSource},
		},
		&core.FileDesc{
			Path: path.Join("infra", "config.go"),
			Data: func() ([]byte, error) {
				var buf bytes.Buffer

				templates.WriteGoBackendInfraConfig(&buf, projectName)
				return buf.Bytes(), nil
			},
			Transforms: []core.Transform{transform.GoFormatSource},
		},
		&core.FileDesc{
			Path: path.Join("infra", "wire_set.go"),
			Data: func() ([]byte, error) {
				var buf bytes.Buffer

				templates.WriteGoBackendInfraWireSet(&buf)
				return buf.Bytes(), nil
			},
			Transforms: []core.Transform{transform.GoFormatSource},
		},
		&core.FileDesc{
			Path: path.Join("infra", "providers.go"),
			Data: func() ([]byte, error) {
				var buf bytes.Buffer

				templates.WriteGoBackendInfraProviders(&buf, projectName, authorName)
				return buf.Bytes(), nil
			},
			Transforms: []core.Transform{transform.GoFormatSource},
		},
		&core.FileDesc{
			Path: path.Join("web", "routes.go"),
			Data: func() ([]byte, error) {
				var buf bytes.Buffer

				templates.WriteGoBackendServerRoutes(&buf, projectName, authorName)
				return buf.Bytes(), nil
			},
			Transforms: []core.Transform{transform.GoFormatSource},
		},
		&core.FileDesc{
			Path: path.Join("web", "server.go"),
			Data: func() ([]byte, error) {
				var buf bytes.Buffer

				templates.WriteGoBackendWebServer(&buf, projectName, authorName)
				return buf.Bytes(), nil
			},
			Transforms: []core.Transform{transform.GoFormatSource},
		},
		&core.FileDesc{
			Path: path.Join("web", "wire.go"),
			Data: func() ([]byte, error) {
				var buf bytes.Buffer

				templates.WriteGoBackendWebWire(&buf, projectName, authorName)
				return buf.Bytes(), nil
			},
			Transforms: []core.Transform{transform.GoFormatSource},
		},
		&core.FileDesc{
			Path: "main.go",
			Data: func() ([]byte, error) {
				var buf bytes.Buffer

				templates.WriteGoBackendMain(&buf, projectName, authorName)
				return buf.Bytes(), nil
			},
			Transforms: []core.Transform{transform.GoFormatSource},
		},
		&core.FileDesc{
			Path: "go.mod",
			Data: func() ([]byte, error) {
				var buf bytes.Buffer

				templates.WriteGoBackendMod(&buf, projectName, authorName, config.GoVersion)
				return buf.Bytes(), nil
			},
		},
		&core.FileDesc{
			Path: ".vscode/settings.json",
			Data: func() ([]byte, error) {
				var buf bytes.Buffer

				templates.WriteGoBackendVscodeSettings(&buf)
				return buf.Bytes(), nil
			},
		},
		&core.FileDesc{
			Path: "Dockerfile",
			Data: func() ([]byte, error) {
				var buf bytes.Buffer

				templates.WriteGoBackendDockerfile(&buf)
				return buf.Bytes(), nil
			},
		},
		&core.FileDesc{
			Path: "Makefile",
			Data: func() ([]byte, error) {
				var buf bytes.Buffer

				templates.WriteGoMakefile(&buf)
				return buf.Bytes(), nil
			},
		},
		&core.FileDesc{
			Path: "air.toml",
			Data: func() ([]byte, error) {
				var buf bytes.Buffer

				templates.WriteGoAirConfig(&buf)
				return buf.Bytes(), nil
			},
		},
		&core.FileDesc{
			Path: "config.yml",
			Data: func() ([]byte, error) {
				var buf bytes.Buffer

				templates.WriteGoBackendConfigYml(&buf, projectName)
				return buf.Bytes(), nil
			},
		},
		&core.FileDesc{
			Path: ".dockerignore",
			Data: func() ([]byte, error) {
				var buf bytes.Buffer

				templates.WriteGoDockerignore(&buf)
				return buf.Bytes(), nil
			},
		},
	).AddCommand(
		&core.ShellCommand{
			Bin:  "go",
			Args: []string{"mod", "tidy"},
			Cwd:  path.Join("./" + projectName),
		},
		&core.ShellCommand{
			Bin:  "wire",
			Args: []string{"./server"},
			Cwd:  path.Join("./" + projectName),
		},
	).Run()
}
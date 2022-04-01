package cmd

import (
	"bytes"
	"os"
	"path"

	"github.com/pot-code/web-cli/pkg/task"
	"github.com/pot-code/web-cli/pkg/transformation"
	"github.com/pot-code/web-cli/pkg/util"
	"github.com/pot-code/web-cli/templates"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type GoWebConfig struct {
	GenType     string `flag:"type" alias:"t" usage:"backend type" validate:"required,oneof=go"`
	ProjectName string `arg:"0" alias:"project_name" validate:"required,var"`
	AuthorName  string `flag:"author" alias:"a" usage:"author name for the app" validate:"required,var"`
	GoVersion   string `flag:"version" alias:"v" usage:"go compiler version" validate:"required,version"`
}

var GoWebCmd = util.NewCliCommand("web", "generate golang web project",
	&GoWebConfig{
		GenType:   "go",
		GoVersion: "1.16",
	},
	util.WithAlias([]string{"w"}),
	util.WithArgUsage("project_name"),
).AddFeature(new(GenWebProject)).ExportCommand()

type GenWebProject struct{}

var _ util.CommandFeature = &GenWebProject{}

func (gwp *GenWebProject) Cond(c *cli.Context) bool {
	return c.String("type") == "go"
}

func (gwp *GenWebProject) Handle(c *cli.Context, cfg interface{}) error {
	config := cfg.(*GoWebConfig)
	projectName := config.ProjectName
	authorName := config.AuthorName

	_, err := os.Stat(projectName)
	if err == nil {
		log.Infof("[skipped]'%s' already exists", projectName)
		return nil
	}

	return util.NewTaskComposer(projectName).AddFile(
		[]*task.FileDesc{
			{
				Path: path.Join("cmd", "web", "main.go"),
				Source: func(buf *bytes.Buffer) error {
					templates.WriteGoServerCmdWebMain(buf, projectName, authorName)
					return nil
				},
				Transforms: []task.Pipeline{transformation.GoFormatSource},
			},
			{
				Path: path.Join("config", "config.go"),
				Source: func(buf *bytes.Buffer) error {
					templates.WriteGoServerConfig(buf)
					return nil
				},
				Transforms: []task.Pipeline{transformation.GoFormatSource},
			},
			{
				Path: path.Join("web", "wire.go"),
				Source: func(buf *bytes.Buffer) error {
					templates.WriteGoServerWebWire(buf, projectName, authorName)
					return nil
				},
				Transforms: []task.Pipeline{transformation.GoFormatSource},
			},
			{
				Path: path.Join("web", "server.go"),
				Source: func(buf *bytes.Buffer) error {
					templates.WriteGoServerWebServer(buf, projectName, authorName)
					return nil
				},
				Transforms: []task.Pipeline{transformation.GoFormatSource},
			},
			{
				Path: path.Join("web", "router.go"),
				Source: func(buf *bytes.Buffer) error {
					templates.WriteGoServerWebRouter(buf)
					return nil
				},
				Transforms: []task.Pipeline{transformation.GoFormatSource},
			},
			{
				Path: "tools.go",
				Source: func(buf *bytes.Buffer) error {
					templates.WriteGoServerTools(buf)
					return nil
				},
				Transforms: []task.Pipeline{transformation.GoFormatSource},
			},
			{
				Path: "go.mod",
				Source: func(buf *bytes.Buffer) error {
					templates.WriteGoMod(buf, projectName, authorName, config.GoVersion)
					return nil
				},
			},
			{
				Path: ".vscode/settings.json",
				Source: func(buf *bytes.Buffer) error {
					templates.WriteGoServerVscodeSettings(buf)
					return nil
				},
			},
			{
				Path: "Dockerfile",
				Source: func(buf *bytes.Buffer) error {
					templates.WriteGoServerDockerfile(buf)
					return nil
				},
			},
			{
				Path: "Makefile",
				Source: func(buf *bytes.Buffer) error {
					templates.WriteGoServerMakefile(buf)
					return nil
				},
			},
			{
				Path: "air.toml",
				Source: func(buf *bytes.Buffer) error {
					templates.WriteGoAirConfig(buf)
					return nil
				},
			},
			{
				Path: "config.yml",
				Source: func(buf *bytes.Buffer) error {
					templates.WriteGoServerConfigYml(buf, projectName)
					return nil
				},
			},
			{
				Path: ".dockerignore",
				Source: func(buf *bytes.Buffer) error {
					templates.WriteGoDockerignore(buf)
					return nil
				},
			},
		}...).AddCommand(
		&task.ShellCommand{
			Bin:  "go",
			Args: []string{"mod", "tidy"},
			Cwd:  path.Join("./" + projectName),
		},
		&task.ShellCommand{
			Bin:  "wire",
			Args: []string{"./web"},
			Cwd:  path.Join("./" + projectName),
		},
	).Run()
}

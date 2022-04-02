package cmd

import (
	"bytes"
	"path"

	"github.com/pot-code/web-cli/internal/command"
	"github.com/pot-code/web-cli/internal/task"
	"github.com/pot-code/web-cli/internal/transformer"
	"github.com/pot-code/web-cli/internal/util"
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

var GoWebCmd = command.NewCliCommand("web", "generate golang web project",
	&GoWebConfig{
		GenType:   "go",
		GoVersion: "1.16",
	},
	command.WithAlias([]string{"w"}),
	command.WithArgUsage("project_name"),
).AddFeature(new(GenWebProject)).ExportCommand()

type GenWebProject struct{}

var _ command.CommandFeature = &GenWebProject{}

func (gwp *GenWebProject) Cond(c *cli.Context) bool {
	return c.String("type") == "go"
}

func (gwp *GenWebProject) Handle(c *cli.Context, cfg interface{}) error {
	config := cfg.(*GoWebConfig)
	projectName := config.ProjectName
	authorName := config.AuthorName

	if util.IsFileExist(projectName) {
		log.Infof("folder '%s' already exists", projectName)
		return nil
	}

	return task.NewSequentialExecutor(
		task.NewParallelExecutor(
			task.BatchFileRequest(
				task.NewFileRequestTree(projectName).
					AddNode( // root/
						&task.FileRequest{
							Name: "tools.go",
							Data: bytes.NewBufferString(templates.GoServerTools()),
						},
					).
					Branch("config").
					AddNode( // root/config
						&task.FileRequest{
							Name: "config.go",
							Data: bytes.NewBufferString(templates.GoServerConfig()),
						},
					).Up().
					Branch("web").
					AddNode( // root/web
						&task.FileRequest{
							Name: "wire.go",
							Data: bytes.NewBufferString(templates.GoServerWebWire(projectName, authorName)),
						},
						&task.FileRequest{
							Name: "server.go",
							Data: bytes.NewBufferString(templates.GoServerWebServer(projectName, authorName)),
						},
						&task.FileRequest{
							Name: "router.go",
							Data: bytes.NewBufferString(templates.GoServerWebRouter()),
						},
					).Up().
					Branch("cmd").Branch("web").
					AddNode( // root/cmd/web
						&task.FileRequest{
							Name: "main.go",
							Data: bytes.NewBufferString(templates.GoServerCmdWebMain(projectName, authorName)),
						},
					).
					Flatten(),
				transformer.GoFormatSource(),
			)...,
		),
		task.NewParallelExecutor(
			task.BatchFileRequest(
				task.NewFileRequestTree(projectName).
					AddNode(
						[]*task.FileRequest{
							{
								Name: "go.mod",
								Data: bytes.NewBufferString(templates.GoMod(projectName, authorName, config.GoVersion)),
							},
							{
								Name: "Dockerfile",
								Data: bytes.NewBufferString(templates.GoServerDockerfile()),
							},
							{
								Name: "Makefile",
								Data: bytes.NewBufferString(templates.GoServerMakefile()),
							},
							{
								Name: "air.toml",
								Data: bytes.NewBufferString(templates.GoAirConfig()),
							},
							{
								Name: "config.yml",
								Data: bytes.NewBufferString(templates.GoServerConfigYml(projectName)),
							},
							{
								Name: ".dockerignore",
								Data: bytes.NewBufferString(templates.GoDockerignore()),
							},
						}...,
					).
					Branch(".vscode").
					AddNode(
						&task.FileRequest{
							Name: "settings.json",
							Data: bytes.NewBufferString(templates.GoServerVscodeSettings()),
						},
					).
					Flatten(),
			)...,
		),
		task.NewShellCmdExecutor(
			&task.ShellCommand{
				Bin:  "go",
				Args: []string{"mod", "tidy"},
				Cwd:  path.Join("./" + projectName),
			},
		),
		task.NewShellCmdExecutor(
			&task.ShellCommand{
				Bin:  "wire",
				Args: []string{"./web"},
				Cwd:  path.Join("./" + projectName),
			},
		),
	).Run()
}

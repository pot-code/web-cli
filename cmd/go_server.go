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
	ProjectName string `arg:"0" alias:"project_name" validate:"required,var"`
	AuthorName  string `flag:"author" alias:"a" usage:"author name for the app" validate:"required,var"`
	GoVersion   string `flag:"version" alias:"v" usage:"go compiler version" validate:"required,version"`
}

var GoWebCmd = command.NewCliCommand("web", "generate golang web project",
	&GoWebConfig{
		GoVersion: "1.16",
	},
	command.WithAlias([]string{"w"}),
	command.WithArgUsage("project_name"),
).AddFeature(GenWebProject).ExportCommand()

var GenWebProject = util.NoCondFeature(func(c *cli.Context, cfg interface{}) error {
	config := cfg.(*GoWebConfig)
	projectName := config.ProjectName
	authorName := config.AuthorName

	if util.IsFileExist(projectName) {
		log.Infof("folder '%s' already exists", projectName)
		return nil
	}

	return task.NewSequentialExecutor(
		task.NewParallelExecutor(
			task.BatchFileTransformation(
				task.NewFileGenerationTree(projectName).
					AddNodes( // root/
						&task.FileGenerator{
							Name: "tools.go",
							Data: bytes.NewBufferString(templates.GoServerTools()),
						},
					).
					Branch("config").
					AddNodes( // root/config
						&task.FileGenerator{
							Name: "config.go",
							Data: bytes.NewBufferString(templates.GoServerConfig()),
						},
					).Up().
					Branch("web").
					AddNodes( // root/web
						&task.FileGenerator{
							Name: "wire.go",
							Data: bytes.NewBufferString(templates.GoServerWebWire(projectName, authorName)),
						},
						&task.FileGenerator{
							Name: "server.go",
							Data: bytes.NewBufferString(templates.GoServerWebServer(projectName, authorName)),
						},
						&task.FileGenerator{
							Name: "router.go",
							Data: bytes.NewBufferString(templates.GoServerWebRouter()),
						},
					).Up().
					Branch("cmd").Branch("web").
					AddNodes( // root/cmd/web
						&task.FileGenerator{
							Name: "main.go",
							Data: bytes.NewBufferString(templates.GoServerCmdWebMain(projectName, authorName)),
						},
					).
					Flatten(),
				transformer.GoFormatSource(),
			)...,
		),
		task.NewParallelExecutor(
			task.BatchFileTransformation(
				task.NewFileGenerationTree(projectName).
					AddNodes(
						[]*task.FileGenerator{
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
					AddNodes(
						&task.FileGenerator{
							Name: "settings.json",
							Data: bytes.NewBufferString(templates.GoServerVscodeSettings()),
						},
					).
					Flatten(),
			)...,
		),
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
})

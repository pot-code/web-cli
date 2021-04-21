package generate

import (
	"bytes"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/pot-code/web-cli/core"
	"github.com/pot-code/web-cli/templates"
	"github.com/pot-code/web-cli/util"
	log "github.com/sirupsen/logrus"
)

type golangApiGenerator struct {
	config   *GolangApiConfig
	runner   core.Generator
	composer *util.TaskComposer
}

type GolangApiConfig struct {
	PackageName string // go pkg name
	Project     string
	Author      string
	Model       string
	Root        string // path root under which to generate ap
}

var _ core.Generator = golangApiGenerator{}

// NewGolangApiGenerator create go api producer
//
// project structure:
//
// <root>
// 	|-<package name>
// 		|-transport
// 			|-http.go
// 		repo.go
// 		service.go
// 		model.go
func NewGolangApiGenerator(config *GolangApiConfig) *golangApiGenerator {
	composer := util.NewTaskComposer(
		path.Join(config.Root, config.PackageName),
		&core.FileDesc{
			Path: "transport/http.go",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoApiHttp(&buf, config.Project, config.Author, config.PackageName, config.Model)
				return buf.Bytes()
			},
		},
		&core.FileDesc{
			Path: "model.go",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoApiModel(&buf, config.PackageName, config.Model)
				return buf.Bytes()
			},
		},
		&core.FileDesc{
			Path: "repo.go",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoApiRepo(&buf, config.PackageName, config.Model)
				return buf.Bytes()
			},
		},
		&core.FileDesc{
			Path: "service.go",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoApiService(&buf, config.PackageName, config.Model)
				return buf.Bytes()
			},
		},
	)
	return &golangApiGenerator{config: config, composer: composer}
}

func (gag golangApiGenerator) Run() error {
	log.Debugf("generation tree:\n%s", gag.composer.GetGenerationTree())
	dir := path.Join(gag.config.Root, gag.config.PackageName)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		runner := gag.composer.MakeRunner()
		gag.runner = runner
		return errors.Wrap(runner.Run(), "failed to generate go api")
	}
	if err == nil {
		log.Infof("[skipped]'%s' already exists", dir)
	}
	return errors.Wrap(err, "failed to generate go api")
}

func (gag golangApiGenerator) Cleanup() error {
	if gag.runner != nil {
		gag.runner.Cleanup()
	}

	root := gag.config.PackageName
	log.Debugf("removing folder '%s'", root)
	err := os.RemoveAll(root)
	if err != nil {
		log.WithFields(log.Fields{"error": err.Error(), "folder": root}).Debug("[cleanup]failed to cleanup")
	}
	return errors.WithStack(err)
}

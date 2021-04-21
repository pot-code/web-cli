package generate

import (
	"bytes"
	"os"

	"github.com/pkg/errors"
	"github.com/pot-code/web-cli/core"
	"github.com/pot-code/web-cli/templates"
	"github.com/pot-code/web-cli/util"
	log "github.com/sirupsen/logrus"
)

type golangBackendGenerator struct {
	config *GolangBackendConfig
	runner core.Generator
}

type GolangBackendConfig struct {
	ProjectName string // project name
	Author      string // project author name
	Version     string // version number
}

var _ core.Generator = golangBackendGenerator{}

// NewGolangBackendGenerator create go backend generator
//
// project structure:
//
// <name>
// 		|-cmd
// 			|-main.go
// 		|-config
// 			|-def.go
// 		|-api
// 			|-api.go
// 		go.mod
func NewGolangBackendGenerator(config *GolangBackendConfig) *golangBackendGenerator {
	composer := util.NewTaskComposer(
		config.ProjectName,
		&core.FileDesc{
			Path: "api/api.go",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoBackendApi(&buf, config.ProjectName, config.Author)
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
			Path: "cmd/main.go",
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
	return &golangBackendGenerator{config: config, runner: composer}
}

func (gbg golangBackendGenerator) Run() error {
	_, err := os.Stat(gbg.config.ProjectName)
	if err == nil {
		log.Infof("[skipped]'%s' already exists", gbg.config.ProjectName)
		return nil
	}
	if os.IsNotExist(err) {
		err = gbg.runner.Run()
	}
	return errors.Wrap(err, "failed to generate go backend")
}

func (gbg golangBackendGenerator) Cleanup() error {
	root := gbg.config.ProjectName
	log.Debugf("removing folder '%s'", root)
	err := os.RemoveAll(root)
	if err != nil {
		log.WithFields(log.Fields{"error": err.Error(), "folder": root}).Debug("[cleanup]failed to cleanup")
	}
	return errors.WithStack(err)
}

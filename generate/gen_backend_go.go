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
	gen    core.Generator
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
	recipe := util.NewGenerationRecipe(
		config.ProjectName,
		&util.GenerationMaterial{
			Path: "api/api.go",
			Provider: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoBackendApi(&buf, config.ProjectName, config.Author)
				return buf.Bytes()
			},
		},
		&util.GenerationMaterial{
			Path: "config/def.go",
			Provider: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoBackendConfig(&buf, config.ProjectName)
				return buf.Bytes()
			},
		},
		&util.GenerationMaterial{
			Path: "cmd/main.go",
			Provider: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoBackendMain(&buf, config.ProjectName, config.Author)
				return buf.Bytes()
			},
		},
		&util.GenerationMaterial{
			Path: "go.mod",
			Provider: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoBackendMod(&buf, config.ProjectName, config.Author, config.Version)
				return buf.Bytes()
			},
		},
	)
	log.Debugf("generation tree:\n%s", recipe.GetGenerationTree())
	return &golangBackendGenerator{config, recipe.MakeGenerator()}
}

func (gbg golangBackendGenerator) Gen() error {
	_, err := os.Stat(gbg.config.ProjectName)
	if os.IsNotExist(err) {
		return errors.Wrap(gbg.gen.Gen(), "failed to generate go backend")
	}
	if err == nil {
		log.Infof("[skipped]'%s' already exists", gbg.config.ProjectName)
	}
	return errors.Wrap(err, "failed to generate go backend")
}

func (gbg golangBackendGenerator) Cleanup() error {
	gbg.gen.Cleanup()

	root := gbg.config.ProjectName
	log.Debugf("removing folder '%s'", root)
	err := os.RemoveAll(root)
	if err != nil {
		log.WithFields(log.Fields{"error": err.Error(), "folder": root}).Debug("[cleanup]failed to cleanup")
	}
	return errors.WithStack(err)
}

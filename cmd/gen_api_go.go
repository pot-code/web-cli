package cmd

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
	config *GenApiConfig
	gen    core.Generator
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
func NewGolangApiGenerator(config *GenApiConfig) *golangApiGenerator {
	recipe := util.NewGenerationRecipe(
		path.Join(config.Root, config.PackageName),
		&util.GenerationMaterial{
			Path: "transport/http.go",
			Provider: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoApiHttp(&buf, config.Project, config.Author, config.PackageName, config.Model)
				return buf.Bytes()
			},
		},
		&util.GenerationMaterial{
			Path: "model.go",
			Provider: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoApiModel(&buf, config.PackageName, config.Model)
				return buf.Bytes()
			},
		},
		&util.GenerationMaterial{
			Path: "repo.go",
			Provider: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoApiRepo(&buf, config.PackageName, config.Model)
				return buf.Bytes()
			},
		},
		&util.GenerationMaterial{
			Path: "service.go",
			Provider: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoApiService(&buf, config.PackageName, config.Model)
				return buf.Bytes()
			},
		},
	)
	log.Debugf("generation tree:\n%s", recipe.GetGenerationTree())
	return &golangApiGenerator{config, recipe.MakeGenerator()}
}

func (gag golangApiGenerator) Gen() error {
	dir := path.Join(gag.config.Root, gag.config.PackageName)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return errors.Wrap(gag.gen.Gen(), "failed to generate go api")
	}
	if err == nil {
		log.Infof("[skipped]'%s' already exists", dir)
	}
	return errors.Wrap(err, "failed to generate go api")
}

func (gag golangApiGenerator) Cleanup() error {
	gag.gen.Cleanup()

	root := gag.config.PackageName
	log.Debugf("removing folder '%s'", root)
	err := os.RemoveAll(root)
	if err != nil {
		log.WithFields(log.Fields{"error": err.Error(), "folder": root}).Debug("[cleanup]failed to cleanup")
	}
	return errors.WithStack(err)
}

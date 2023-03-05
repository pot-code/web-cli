package provider

import (
	"embed"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
)

var templateFs embed.FS

func InitTemplateFS(fs embed.FS) {
	templateFs = fs
}

type EmbedFileProvider struct {
	p string
}

func NewEmbedFileProvider(p string) *EmbedFileProvider {
	return &EmbedFileProvider{p}
}

func (e *EmbedFileProvider) Get() (io.Reader, error) {
	log.WithFields(log.Fields{
		"path":     e.p,
		"provider": "EmbedFileProvider",
	}).Debug("open file")

	fd, err := templateFs.Open(e.p)
	if err != nil {
		return nil, fmt.Errorf("open file at %s: %w", e.p, err)
	}
	return fd, nil
}

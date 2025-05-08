package provider

import (
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog/log"
)

type LocalFileProvider struct {
	p string
}

func NewLocalFileProvider(p string) *LocalFileProvider {
	return &LocalFileProvider{p}
}

func (p *LocalFileProvider) Get() (io.Reader, error) {
	log.Debug().
		Str("path", p.p).
		Str("provider", "LocalFileProvider").
		Msg("open file")

	fd, err := os.Open(p.p)
	if err != nil {
		return nil, fmt.Errorf("open file at %s: %w", p.p, err)
	}
	return fd, nil
}

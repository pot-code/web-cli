package provider

import (
	"fmt"
	"io"
	"os"
)

type LocalFileProvider struct {
	p string
}

func NewLocalFileProvider(p string) *LocalFileProvider {
	return &LocalFileProvider{p}
}

func (p *LocalFileProvider) Get() (io.Reader, error) {
	fd, err := os.Open(p.p)
	if err != nil {
		return nil, fmt.Errorf("open file at %s: %w", p.p, err)
	}
	return fd, nil
}

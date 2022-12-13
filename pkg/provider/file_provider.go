package provider

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type LocalFileProvider struct {
	Path string
}

func NewLocalFileProvider(p string) *LocalFileProvider {
	return &LocalFileProvider{Path: p}
}

func (p *LocalFileProvider) Get() (io.Reader, error) {
	fd, err := os.Open(p.Path)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	return fd, nil
}

const ConnectionTimeout = 30 * time.Second

type RemoteFileProvider struct {
	URL string
}

func NewRemoteFileProvider(url string) *RemoteFileProvider {
	return &RemoteFileProvider{URL: url}
}

func (p *RemoteFileProvider) Get() (io.Reader, error) {
	conn := http.Client{
		Timeout: ConnectionTimeout,
	}
	res, err := conn.Get(p.URL)
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}

	return res.Body, nil
}

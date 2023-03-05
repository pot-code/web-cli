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
		return nil, fmt.Errorf("open file at %s: %w", p.Path, err)
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
		return nil, fmt.Errorf("get file from %s: %w", p.URL, err)
	}
	return res.Body, nil
}

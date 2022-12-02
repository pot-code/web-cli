package task

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type LocalTemplateProvider struct {
	Path string
}

func NewLocalTemplateProvider(p string) *LocalTemplateProvider {
	return &LocalTemplateProvider{Path: p}
}

func (p *LocalTemplateProvider) Get() (io.ReadCloser, error) {
	fd, err := os.Open(p.Path)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	return fd, nil
}

const ConnectionTimeout = 30 * time.Second

type RemoteTemplateProvider struct {
	URL string
}

func NewRemoteTemplateProvider(url string) *RemoteTemplateProvider {
	return &RemoteTemplateProvider{URL: url}
}

func (p *RemoteTemplateProvider) Get() (io.ReadCloser, error) {
	conn := http.Client{
		Timeout: ConnectionTimeout,
	}
	res, err := conn.Get(p.URL)
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}

	return res.Body, nil
}

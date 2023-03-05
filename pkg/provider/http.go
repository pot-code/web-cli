package provider

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

const ConnectionTimeout = 30 * time.Second

type RemoteFileProvider struct {
	url string
}

func NewRemoteFileProvider(url string) *RemoteFileProvider {
	return &RemoteFileProvider{url}
}

func (p *RemoteFileProvider) Get() (io.Reader, error) {
	conn := http.Client{
		Timeout: ConnectionTimeout,
	}
	res, err := conn.Get(p.url)
	if err != nil {
		return nil, fmt.Errorf("get file from %s: %w", p.url, err)
	}
	return res.Body, nil
}

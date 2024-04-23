package provider

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

const ConnectionTimeout = 30 * time.Second

type RemoteFileProvider struct {
	url string
}

func NewRemoteFileProvider(url string) *RemoteFileProvider {
	return &RemoteFileProvider{url}
}

func (p *RemoteFileProvider) Get() (io.Reader, error) {
	log.Debug().
		Str("path", p.url).
		Str("provider", "RemoteFileProvider").
		Dur("timeout", ConnectionTimeout).
		Msg("fetch file")

	conn := http.Client{
		Timeout: ConnectionTimeout,
	}
	res, err := conn.Get(p.url)
	if err != nil {
		return nil, fmt.Errorf("get file from %s: %w", p.url, err)
	}
	return res.Body, nil
}

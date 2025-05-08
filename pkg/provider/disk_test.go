package provider

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalProvider_Get(t *testing.T) {
	p := NewLocalFileProvider("./__fixture__/test.tmpl")

	rc, err := p.Get()
	assert.Nil(t, err)

	c, err := io.ReadAll(rc)
	assert.Nil(t, err)

	assert.Equal(t, "hello {{.Name}}", string(c))
}

func TestRemoteProvider_Get(t *testing.T) {
	p := NewRemoteFileProvider("https://raw.githubusercontent.com/pot-code/react-template/master/.nvmrc")

	rc, err := p.Get()
	assert.Nil(t, err)

	c, err := io.ReadAll(rc)
	assert.Nil(t, err)

	assert.Equal(t, "v18.12.1\n", string(c))
}

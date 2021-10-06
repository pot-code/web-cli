package transform

import (
	"bytes"
	"go/format"

	"github.com/pkg/errors"
)

func GoFormatSource(source *bytes.Buffer) (*bytes.Buffer, error) {
	dest := new(bytes.Buffer)
	fs, err := format.Source(source.Bytes())
	dest.Write(fs)
	return dest, errors.Wrap(err, "failed to format")
}

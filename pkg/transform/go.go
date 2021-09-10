package transform

import (
	"go/format"

	"github.com/pkg/errors"
)

func GoFormatSource(source []byte) ([]byte, error) {
	fs, err := format.Source(source)
	return fs, errors.Wrap(err, "failed to format")
}

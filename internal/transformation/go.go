package transformation

import (
	"bytes"
	"go/format"

	"github.com/pkg/errors"
)

func GoFormatSource(src *bytes.Buffer, dst *bytes.Buffer) error {
	fs, err := format.Source(src.Bytes())
	dst.Write(fs)
	return errors.Wrap(err, "failed to format")
}

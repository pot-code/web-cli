package transformer

import (
	"bytes"
	"go/format"

	"github.com/pkg/errors"
)

func GoFormatSource() Transformer {
	return func(next TransformerFunc) TransformerFunc {
		return func(data *bytes.Buffer) error {
			result, err := format.Source(data.Bytes())
			if err != nil {
				return errors.Wrap(err, "failed to format Go source")
			}
			return next(bytes.NewBuffer(result))
		}
	}
}

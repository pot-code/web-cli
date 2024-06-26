package pipeline

import (
	"fmt"
	"go/format"
	"io"

	"github.com/rs/zerolog/log"
)

func GoFormatSource(src io.Reader, dest io.Writer) error {
	d, err := io.ReadAll(src)
	if err != nil {
		return fmt.Errorf("format go source: read from src: %w", err)
	}

	result, err := format.Source(d)
	if err != nil {
		return fmt.Errorf("format go source: format: %w", err)
	}
	n, err := dest.Write(result)
	if err != nil {
		return fmt.Errorf("format go source: write formatted code: %w", err)
	}
	log.Debug().Int("bytes", n).Msg("write formatted code")
	return nil
}

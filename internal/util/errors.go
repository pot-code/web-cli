package util

import (
	"fmt"

	"github.com/pkg/errors"
)

type StackTracer interface {
	StackTrace() errors.StackTrace
}

func GetVerboseStackTrace(depth int, st StackTracer) string {
	frames := st.StackTrace()
	if depth > 0 {
		frames = frames[:depth]
	}
	return fmt.Sprintf("%+v", frames)
}
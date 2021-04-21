package util

import (
	"fmt"

	"github.com/pkg/errors"
)

// CommandError command misuse error
type CommandError struct {
	cmd string
	err error
}

func NewCommandError(cmd string, err error) *CommandError {
	return &CommandError{
		cmd, err,
	}
}

func (ce CommandError) Error() string {
	return ce.err.Error()
}

func (ce CommandError) Command() string {
	return ce.cmd
}

func (ce CommandError) Unwrap() error {
	return ce.err
}

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

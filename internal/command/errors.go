package command

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

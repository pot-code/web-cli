package command

import (
	"github.com/urfave/cli/v2"
)

type CommandHandler[T any] interface {
	Handle(c *cli.Context, cfg T) error
}

type InlineHandler[T any] func(c *cli.Context, cfg T) error

func (s InlineHandler[T]) Handle(c *cli.Context, cfg T) error {
	return s(c, cfg)
}

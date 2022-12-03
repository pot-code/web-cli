package command

import (
	"github.com/urfave/cli/v2"
)

type CommandHandler interface {
	Handle(c *cli.Context, cfg interface{}) error
}

type InlineHandler func(c *cli.Context, cfg interface{}) error

func (s InlineHandler) Handle(c *cli.Context, cfg interface{}) error {
	return s(c, cfg)
}

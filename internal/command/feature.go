package command

import (
	"github.com/urfave/cli/v2"
)

type CommandFeature interface {
	Cond(c *cli.Context) bool
	Handle(c *cli.Context, cfg interface{}) error
}

type NoCondFeature func(c *cli.Context, cfg interface{}) error

func (s NoCondFeature) Handle(c *cli.Context, cfg interface{}) error {
	return s(c, cfg)
}

func (s NoCondFeature) Cond(c *cli.Context) bool {
	return true
}

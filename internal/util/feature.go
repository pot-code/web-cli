package util

import (
	"github.com/urfave/cli/v2"
)

type NoCondFeature func(c *cli.Context, cfg interface{}) error

func (s NoCondFeature) Handle(c *cli.Context, cfg interface{}) error {
	return s(c, cfg)
}

func (s NoCondFeature) Cond(c *cli.Context) bool {
	return true
}

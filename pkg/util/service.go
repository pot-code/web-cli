package util

import "github.com/urfave/cli/v2"

type NoCondFunctionService func(c *cli.Context, cfg interface{}) error

func (s NoCondFunctionService) Handle(c *cli.Context, cfg interface{}) error {
	return s(c, cfg)
}

func (s NoCondFunctionService) Cond(c *cli.Context) bool {
	return true
}

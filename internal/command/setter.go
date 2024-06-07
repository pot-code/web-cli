package command

import (
	"reflect"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

type configSetter struct {
	flags []*flagField
	args  []*argField
	avg   *argValueParser
}

func newConfigSetter(flags []*flagField, args []*argField) *configSetter {
	return &configSetter{flags: flags, args: args, avg: newArgValueParser()}
}

func (c *configSetter) setFromContext(ctx *cli.Context, config interface{}) error {
	if err := validateConfig(config); err != nil {
		return nil
	}

	rv := reflect.ValueOf(config).Elem()

	for _, f := range c.flags {
		rv.FieldByName(f.fieldName).Set(reflect.ValueOf(ctx.Value(f.flagName)))
	}
	for _, f := range c.args {
		v, err := c.avg.parse(f, ctx.Args().Get(f.position))
		if err != nil {
			log.Err(err).Str("field", f.fieldName).Int("position", f.position).Msg("parse arg value")
			return err
		}
		rv.FieldByName(f.fieldName).Set(reflect.ValueOf(v))
	}
	return nil
}

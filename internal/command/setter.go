package command

import (
	"reflect"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

type configReader struct {
	flags []*flagField
	args  []*argField
	avg   *argValueParser
}

func newConfigReader(flags []*flagField, args []*argField) *configReader {
	return &configReader{flags: flags, args: args, avg: newArgValueParser()}
}

func (c *configReader) readFromCliContext(ctx *cli.Context, receiver any) error {
	if err := validateConfig(receiver); err != nil {
		return nil
	}

	rv := reflect.ValueOf(receiver).Elem()

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

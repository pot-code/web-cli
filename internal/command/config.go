package command

import (
	"errors"
	"reflect"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

type configParser struct {
	flags []*flagField
	args  []*argField
	avg   *argValueParser
}

func newConfigParser(flags []*flagField, args []*argField) *configParser {
	return &configParser{flags: flags, args: args, avg: newArgValueParser()}
}

func (c *configParser) parseFromCliContext(ctx *cli.Context, receiver any) error {
	if err := validateConfig(receiver); err != nil {
		return nil
	}

	rv := reflect.ValueOf(receiver).Elem()

	for _, f := range c.flags {
		passValue := reflect.ValueOf(ctx.Value(f.flagName))
		if !passValue.IsZero() {
			rv.FieldByName(f.fieldName).Set(reflect.ValueOf(ctx.Value(f.flagName)))
		}
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

func validateConfig(config any) error {
	if config == nil {
		return errors.New("config is nil")
	}

	rv := reflect.ValueOf(config)
	if rv.Kind() != reflect.Ptr && rv.Elem().Kind() != reflect.Struct {
		return errors.New("config must be of pointer of struct type")
	}
	return nil
}

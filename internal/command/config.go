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
	ap    *argValueParser
}

func newConfigParser(flags []*flagField, args []*argField) *configParser {
	return &configParser{flags: flags, args: args, ap: newArgValueParser()}
}

// parseFromCliContext 从 cli 上下文中解析配置，receiver 为指向配置的指针，用于接收解析后的值
func (c *configParser) parseFromCliContext(ctx *cli.Context, receiver any) error {
	if err := validateConfig(receiver); err != nil {
		return err
	}

	configStruct := reflect.ValueOf(receiver).Elem()

	for _, f := range c.flags {
		ctxValue := ctx.Value(f.flagName)
		value := reflect.ValueOf(ctxValue)
		if !value.IsZero() {
			// 如果 cli 传递的值不是字段类型的 zero 值，那么使用 cli 传递的值
			configStruct.FieldByName(f.fieldName).Set(value)
		}
	}

	for _, f := range c.args {
		v, err := c.ap.parse(f, ctx.Args().Get(f.position))
		if err != nil {
			log.Err(err).Str("field", f.fieldName).Int("position", f.position).Msg("parse arg value")
			return err
		}
		configStruct.FieldByName(f.fieldName).Set(reflect.ValueOf(v))
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

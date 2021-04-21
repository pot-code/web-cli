package util

import (
	"reflect"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

// ParseConfig parse and return validation error or parse error
//
// config: receiver
func ParseConfig(c *cli.Context, config interface{}) error {
	t := reflect.TypeOf(config)
	if t.Kind() != reflect.Ptr {
		return errors.New("config must be of pointer type")
	}
	t = t.Elem()

	v := reflect.ValueOf(config)
	v = reflect.Indirect(v)
	for i := v.NumField(); i > 0; i-- {
		f := t.Field(i - 1)

		if arg := f.Tag.Get("arg"); arg != "" {
			index, err := strconv.Atoi(arg)
			if err != nil {
				return errors.Wrapf(err, "failed to convert '%s' to number", arg)
			}

			if f.Type.Kind() == reflect.String {
				v.Field(i - 1).SetString(c.Args().Get(index))
			} else {
				panic("not implemented")
			}
		} else if name := f.Tag.Get("name"); name != "" {
			if f.Type.Kind() == reflect.String {
				v.Field(i - 1).SetString(c.String(name))
			} else if f.Type.Kind() == reflect.Bool {
				v.Field(i - 1).SetBool(c.Bool(name))
			} else {
				panic("not implemented")
			}
		}
	}

	err := validate.Struct(config)
	if v, ok := err.(validator.ValidationErrors); ok {
		msg := v[0].Translate(trans)
		return NewCommandError(c.Command.Name, errors.New(msg))
	}
	return err
}

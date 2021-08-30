package util

import (
	"reflect"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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
	for i := v.NumField() - 1; i >= 0; i-- {
		tf := t.Field(i)
		vf := v.Field(i)

		if !vf.CanSet() {
			log.WithFields(log.Fields{
				"caller": "ParseConfig",
			}).Fatalf("config field [%s] can't be set, maybe it's not exported", tf.Name)
			continue
		}

		if arg := tf.Tag.Get("arg"); arg != "" {
			index, err := strconv.Atoi(arg)
			if err != nil {
				return errors.Wrapf(err, "failed to convert '%s' to number", arg)
			}

			if tf.Type.Kind() == reflect.String {
				vf.SetString(c.Args().Get(index))
			} else {
				panic("not implemented")
			}
		} else if name := tf.Tag.Get("name"); name != "" {
			if tf.Type.Kind() == reflect.String {
				vf.SetString(c.String(name))
			} else if tf.Type.Kind() == reflect.Bool {
				vf.SetBool(c.Bool(name))
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

package command

import (
	"errors"
	"reflect"
)

func validateConfig(config interface{}) error {
	if config == nil {
		return errors.New("config is nil")
	}

	rv := reflect.ValueOf(config)
	if rv.Kind() != reflect.Ptr && rv.Elem().Kind() != reflect.Struct {
		return errors.New("config must be of pointer of struct type")
	}
	return nil
}

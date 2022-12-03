package command

import (
	"fmt"
	"reflect"
)

func walkConfig(config interface{}, v configVisitor) error {
	if config == nil {
		return nil
	}

	if !isPointerType(config) {
		panic("config value must be of pointer type")
	}

	configType := reflect.TypeOf(config).Elem()
	if configType.Kind() != reflect.Struct {
		panic("config must be of struct type")
	}

	var err error
	configValue := reflect.ValueOf(config).Elem()
	for i := configValue.NumField() - 1; i >= 0; i-- {
		ft := configType.Field(i)
		fv := configValue.Field(i)
		if shouldSkipField(ft, fv) {
			continue
		}

		cf := newConfigField(ft, fv)
		switch cf.fieldKind() {
		case reflect.String:
			err = v.visitString(cf)
		case reflect.Int:
			err = v.visitInt(cf)
		case reflect.Bool:
			err = v.visitBoolean(cf)
		case reflect.Slice:
			err = v.visitSlice(cf)
		default:
			panic(fmt.Errorf("unsupported config field kind: %s", cf.fieldType()))
		}

		if err != nil {
			return err
		}
	}
	return err
}

func shouldSkipField(t reflect.StructField, v reflect.Value) bool {
	if t.Tag == "" || !t.IsExported() {
		return true
	}

	noFlag := false
	if t.Tag.Get("flag") == "" {
		noFlag = true
	}

	noArg := false
	if t.Tag.Get("arg") == "" {
		noArg = true
	}

	return noFlag && noArg
}

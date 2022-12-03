package command

import (
	"errors"
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
		cf := newConfigField(configType.Field(i), configValue.Field(i))
		if shouldSkipField(cf) {
			continue
		}

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

func shouldSkipField(f *configField) bool {
	if !f.hasTag() || !f.isExported() {
		return true
	}

	noFlag := false
	if _, err := getFieldFlag(f); errors.Is(err, errEmptyTag) {
		noFlag = true
	}

	noArg := false
	if _, err := getFieldArgPosition(f); errors.Is(err, errEmptyTag) {
		noArg = true
	}

	return noFlag && noArg
}

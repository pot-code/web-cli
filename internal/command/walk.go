package command

import (
	"reflect"
)

type visitor interface {
	accept(f *configField)
}

func walkConfig(config interface{}, v visitor) {
	if config == nil {
		return
	}

	if !isPointerType(config) {
		panic("config must be of pointer type")
	}

	root := newConfigField(
		"",
		reflect.StructField{Name: "", Type: reflect.TypeOf(config).Elem()},
		reflect.ValueOf(config).Elem(),
	)
	walk(root, v)
}

func isPointerType(v interface{}) bool {
	if v == nil {
		return false
	}
	return reflect.TypeOf(v).Kind() == reflect.Ptr
}

func walk(f *configField, v visitor) {
	meta := f.meta
	fieldType := meta.Type
	fieldValue := f.value

	if fieldType.Kind() == reflect.Struct {
		for i := fieldType.NumField() - 1; i >= 0; i-- {
			childField := fieldType.Field(i)
			childValue := fieldValue.Field(i)
			walk(newConfigField(fromParentPrefix(f, fieldType.Name()), childField, childValue), v)
		}
	} else {
		v.accept(f)
	}
}

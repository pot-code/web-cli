package command

import (
	"reflect"
)

func walkConfig(config interface{}, v configVisitor) {
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

func walk(f *configField, v configVisitor) {
	meta := f.structField
	fieldType := meta.Type
	fieldValue := f.value

	if fieldType.Kind() == reflect.Struct {
		for i := fieldType.NumField() - 1; i >= 0; i-- {
			cf := fieldType.Field(i)
			cv := fieldValue.Field(i)
			walk(newConfigField(getQualifiedName(f, fieldType.Name()), cf, cv), v)
		}
	} else {
		v.visit(f)
	}
}

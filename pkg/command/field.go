package command

import (
	"reflect"
)

type configField struct {
	structField reflect.StructField
	value       reflect.Value
}

func newConfigField(structField reflect.StructField, value reflect.Value) *configField {
	return &configField{
		structField: structField,
		value:       value,
	}
}

func (f *configField) fieldType() reflect.Type {
	return f.structField.Type
}

func (f *configField) fieldKind() reflect.Kind {
	return f.structField.Type.Kind()
}

func (f *configField) isExported() bool {
	return f.structField.IsExported()
}

func (f *configField) hasTag() bool {
	return f.structField.Tag != ""
}

func (f *configField) hasDefaultValue() bool {
	return !f.value.IsZero()
}

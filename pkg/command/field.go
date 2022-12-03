package command

import (
	"reflect"
)

type configField struct {
	namespace   string
	structField reflect.StructField
	value       reflect.Value
}

func newConfigField(namespace string, structField reflect.StructField, value reflect.Value) *configField {
	return &configField{
		namespace:   namespace,
		structField: structField,
		value:       value,
	}
}

func (f *configField) qualifiedName() string {
	return getQualifiedName(f, f.structField.Name)
}

func (f *configField) fieldType() reflect.Type {
	return f.structField.Type
}

func (f *configField) fieldKind() reflect.Kind {
	return f.structField.Type.Kind()
}

func (f *configField) typeString() string {
	return f.structField.Type.Kind().String()
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

func getQualifiedName(container *configField, fieldName string) string {
	if container.namespace == "" {
		return fieldName
	}
	return container.namespace + "." + fieldName
}

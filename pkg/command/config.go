package command

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	errEmptyTag = errors.New("empty tag")
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

func (f *configField) getFieldType() reflect.Type {
	return f.structField.Type
}

func (f *configField) getFieldKind() reflect.Kind {
	return f.structField.Type.Kind()
}

func (f *configField) hasDefaultValue() bool {
	return !f.value.IsZero()
}

func (f *configField) getFlag() string {
	return f.getTag("flag")
}

func (f *configField) getAlias() []string {
	t := f.getTag("alias")
	if t == "" {
		return nil
	}
	return strings.Split(t, ",")
}

func (f *configField) getArgPosition() (int, error) {
	i := f.getTag("arg")
	if i == "" {
		return -1, errEmptyTag
	}
	return strconv.Atoi(i)
}

func (f *configField) getUsage() string {
	usage := f.getTag("usage")
	if usage == "" {
		return ""
	}

	choices := f.getChoices()
	if len(choices) == 0 {
		return usage
	}
	return fmt.Sprintf("%s (options: %s)", usage, strings.Join(choices, ", "))
}

func (f *configField) isRequired() bool {
	t := f.getTag("validate")
	if t == "" {
		return false
	}
	return strings.Contains(t, "required")
}

func (f *configField) getChoices() []string {
	vt := f.getTag("validate")
	if vt == "" {
		return nil
	}

	options := strings.Split(vt, ",")
	for _, o := range options {
		if strings.HasPrefix(o, "oneof") {
			return strings.Split(strings.TrimPrefix(o, "oneof="), " ")
		}
	}
	return nil
}

func (f *configField) getTag(t string) string {
	return f.structField.Tag.Get(t)
}

package command

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type configField struct {
	namespace string
	meta      reflect.StructField
	value     reflect.Value
}

func newConfigField(namespace string, meta reflect.StructField, value reflect.Value) *configField {
	return &configField{
		namespace: namespace,
		meta:      meta,
		value:     value,
	}
}

func (f *configField) name() string {
	return fromParentPrefix(f, f.meta.Name)
}

func (f *configField) fieldType() reflect.Type {
	return f.meta.Type
}

func (f *configField) kind() reflect.Kind {
	return f.meta.Type.Kind()
}

func (f *configField) typeString() string {
	return f.meta.Type.Kind().String()
}

func (f *configField) isExported() bool {
	return f.meta.IsExported()
}

func (f *configField) hasTag() bool {
	return f.meta.Tag != ""
}

func (f *configField) hasDefaultValue() bool {
	return !f.value.IsZero()
}

func fromParentPrefix(parent *configField, newPrefix string) string {
	if parent.namespace == "" {
		return newPrefix
	}
	return parent.namespace + "." + newPrefix
}

type help struct {
	base    string
	options []string
}

func (h help) String() string {
	if len(h.options) == 0 {
		return h.base
	}
	return fmt.Sprintf("%s (options: %s)", h.base, strings.Join(h.options, ", "))
}

func getFlag(f *configField) (string, error) {
	tag := f.meta.Tag
	if flag := tag.Get("flag"); flag != "" {
		return flag, nil
	}
	return "", errors.New("empty flag")
}

func getAlias(f *configField) ([]string, error) {
	tag := f.meta.Tag
	if alias := tag.Get("alias"); alias != "" {
		return strings.Split(alias, ","), nil
	}
	return nil, errors.New("empty alias")
}

func getArgPosition(f *configField) (int, error) {
	tag := f.meta.Tag
	if arg := tag.Get("arg"); arg != "" {
		return strconv.Atoi(arg)
	}
	return 0, errors.New("empty arg")
}

func getUsage(f *configField) (string, error) {
	tag := f.meta.Tag
	usage := tag.Get("usage")
	if usage == "" {
		return "", errors.New("empty help")
	}

	h := help{base: usage, options: getOptions(f)}
	return h.String(), nil
}

// isRequired check if field is required
func isRequired(f *configField) bool {
	tag := f.meta.Tag
	if tag == "" {
		return false
	}

	vt := tag.Get("validate")
	if vt == "" {
		return false
	}

	return strings.Contains(vt, "required")
}

// getOptions parse and get oneof items from field tag
func getOptions(f *configField) []string {
	tag := f.meta.Tag
	if tag == "" {
		return nil
	}

	vt := tag.Get("validate")
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

package command

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
)

var (
	errEmptyTag = errors.New("empty tag")
)

var _ configVisitor = (*extractFlagsVisitor)(nil)
var _ configVisitor = (*setConfigVisitor)(nil)

type extractFlagsVisitor struct {
	flags []cli.Flag
}

func newExtractFlagsVisitor() *extractFlagsVisitor {
	return &extractFlagsVisitor{flags: make([]cli.Flag, 0)}
}

func (efv *extractFlagsVisitor) visitSlice(f *configField) error {
	et := f.fieldType().Elem()
	switch et.Kind() {
	case reflect.String:
		efv.visitStringSlice(f)
	default:
		panic(fmt.Errorf("unsupported slice kind '%s'", et.Kind()))
	}
	return nil
}

func (efv *extractFlagsVisitor) visitStringSlice(f *configField) error {
	flag, err := getFieldFlag(f)
	if err != nil {
		return nil
	}

	cf := &cli.StringSliceFlag{
		Name: flag,
	}

	if f.hasDefaultValue() {
		cf.Value = cli.NewStringSlice(reflectValueToStringSlice(f.value)...)
	}

	if isFieldRequired(f) {
		cf.Required = true
	}

	if u, err := getFieldUsage(f); err == nil {
		cf.Usage = u
	}

	if alias, err := getFieldAlias(f); err == nil {
		cf.Aliases = alias
	}

	efv.flags = append(efv.flags, cf)
	return nil
}

func (efv *extractFlagsVisitor) visitString(f *configField) error {
	flag, err := getFieldFlag(f)
	if err != nil {
		return nil
	}

	cf := &cli.StringFlag{
		Name: flag,
	}

	if f.hasDefaultValue() {
		cf.Value = f.value.String()
	}

	if isFieldRequired(f) {
		cf.Required = true
	}

	if u, err := getFieldUsage(f); err == nil {
		cf.Usage = u
	}

	if alias, err := getFieldAlias(f); err == nil {
		cf.Aliases = alias
	}

	efv.flags = append(efv.flags, cf)
	return nil
}

func (efv *extractFlagsVisitor) visitBoolean(f *configField) error {
	flag, err := getFieldFlag(f)
	if err != nil {
		return nil
	}

	cf := &cli.BoolFlag{
		Name: flag,
	}

	if f.hasDefaultValue() {
		cf.Value = f.value.Bool()
	}

	if isFieldRequired(f) {
		cf.Required = true
	}

	if u, err := getFieldUsage(f); err == nil {
		cf.Usage = u
	}

	if alias, err := getFieldAlias(f); err == nil {
		cf.Aliases = alias
	}

	efv.flags = append(efv.flags, cf)
	return nil
}

func (efv *extractFlagsVisitor) visitInt(f *configField) error {
	flag, err := getFieldFlag(f)
	if err != nil {
		return nil
	}

	cf := &cli.IntFlag{
		Name: flag,
	}

	if f.hasDefaultValue() {
		cf.Value = int(f.value.Int())
	}

	if isFieldRequired(f) {
		cf.Required = true
	}

	if u, err := getFieldUsage(f); err == nil {
		cf.Usage = u
	}

	if alias, err := getFieldAlias(f); err == nil {
		cf.Aliases = alias
	}

	efv.flags = append(efv.flags, cf)
	return nil
}

func (efv *extractFlagsVisitor) getFlags() []cli.Flag {
	return efv.flags
}

type setConfigVisitor struct {
	ctx *cli.Context
}

func newSetConfigVisitor(ctx *cli.Context) *setConfigVisitor {
	return &setConfigVisitor{ctx}
}

func (scv *setConfigVisitor) visitSlice(f *configField) error {
	et := f.fieldType().Elem()
	var err error
	switch et.Kind() {
	case reflect.String:
		err = scv.visitStringSlice(f)
	default:
		panic(fmt.Errorf("unsupported slice kind '%s'", et.Kind()))
	}
	return err
}

func (scv *setConfigVisitor) visitStringSlice(f *configField) error {
	var value []string
	ctx := scv.ctx
	if flag, err := getFieldFlag(f); err == nil {
		value = ctx.StringSlice(flag)
	}
	f.value.Set(reflect.ValueOf(value))
	return nil
}

func (scv *setConfigVisitor) visitString(f *configField) error {
	var value string
	ctx := scv.ctx
	if flag, err := getFieldFlag(f); err == nil {
		value = ctx.String(flag)
	} else {
		pos, _ := getFieldArgPosition(f)
		value = ctx.Args().Get(pos)
	}
	f.value.SetString(value)
	return nil
}

func (scv *setConfigVisitor) visitBoolean(f *configField) error {
	flag, _ := getFieldFlag(f)
	value := scv.ctx.Bool(flag)
	f.value.SetBool(value)
	return nil
}

func (scv *setConfigVisitor) visitInt(f *configField) error {
	var value int
	ctx := scv.ctx
	if flag, err := getFieldFlag(f); err == nil {
		value = ctx.Int(flag)
	} else {
		pos, _ := getFieldArgPosition(f)
		av := ctx.Args().Get(pos)
		iv, err := strconv.Atoi(av)
		if err != nil {
			return fmt.Errorf("parse %s to int", av)
		}
		value = iv
	}
	f.value.SetInt(int64(value))
	return nil
}

// isFieldRequired check if field is required
func isFieldRequired(f *configField) bool {
	t, err := getFieldTag(f, "validate")
	if err != nil {
		return false
	}
	return strings.Contains(t, "required")
}

// getFieldChoices parse and get oneof items from field tag
func getFieldChoices(f *configField) []string {
	tag := f.structField.Tag
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

func getFieldFlag(f *configField) (string, error) {
	return getFieldTag(f, "flag")
}

func getFieldAlias(f *configField) ([]string, error) {
	t, err := getFieldTag(f, "alias")
	if err != nil {
		return nil, err
	}
	return strings.Split(t, ","), nil
}

func getFieldArgPosition(f *configField) (int, error) {
	i, err := getFieldTag(f, "arg")
	if err != nil {
		return -1, err
	}

	return strconv.Atoi(i)
}

func getFieldUsage(f *configField) (string, error) {
	usage, err := getFieldTag(f, "usage")
	if err == nil {
		return usage, err
	}

	choices := getFieldChoices(f)
	if len(choices) == 0 {
		return usage, nil
	}
	return fmt.Sprintf("%s (options: %s)", usage, strings.Join(choices, ", ")), nil
}

func getFieldTag(f *configField, t string) (string, error) {
	tag := f.structField.Tag.Get(t)
	if tag == "" {
		return "", errEmptyTag
	}
	return tag, nil
}

func reflectValueToStringSlice(f reflect.Value) []string {
	l := f.Len()
	s := make([]string, l)
	for i := l - 1; i >= 0; i-- {
		s[i] = f.Index(i).String()
	}
	return s
}

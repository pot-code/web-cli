package command

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type configVisitor interface {
	visit(f *configField)
}

type extractFlagsVisitor struct {
	flags []cli.Flag
}

func newExtractFlagsVisitor() *extractFlagsVisitor {
	return &extractFlagsVisitor{flags: make([]cli.Flag, 0)}
}

var _ configVisitor = (*extractFlagsVisitor)(nil)

func (efv *extractFlagsVisitor) visit(f *configField) {
	if !f.isExported() {
		log.Debugf("config field '%s' is not exported", f.qualifiedName())
		return
	}

	if !f.hasTag() {
		log.Debugf("config field '%s' has no tag", f.qualifiedName())
		return
	}

	if _, err := getFieldFlag(f); err != nil {
		log.Debugf("config field '%s' has no flag", f.qualifiedName())
		return
	}

	switch f.fieldKind() {
	case reflect.String:
		efv.visitString(f)
	case reflect.Int:
		efv.visitInt(f)
	case reflect.Bool:
		efv.visitBoolean(f)
	case reflect.Slice:
		efv.visitSlice(f)
	default:
		panic(fmt.Errorf("unsupported kind '%s'", f.fieldKind()))
	}
}

func (efv *extractFlagsVisitor) getFlags() []cli.Flag {
	return efv.flags
}

func (efv *extractFlagsVisitor) visitSlice(f *configField) {
	et := f.fieldType().Elem()
	switch et.Kind() {
	case reflect.String:
		efv.visitStringSlice(f)
	default:
		panic(fmt.Errorf("unsupported slice kind '%s'", et.Kind()))
	}
}

func (efv *extractFlagsVisitor) visitStringSlice(f *configField) {
	flag, _ := getFieldFlag(f)
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

	log.WithFields(log.Fields{
		"flag":      cf.Name,
		"meta_name": f.qualifiedName(),
		"type":      f.typeString(),
		"required":  cf.Required,
		"alias":     cf.Aliases,
		"usage":     cf.Usage,
		"default":   cf.Value,
	}).Debug("append flag")
	efv.flags = append(efv.flags, cf)
}

func (efv *extractFlagsVisitor) visitString(f *configField) {
	flag, _ := getFieldFlag(f)
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

	log.WithFields(log.Fields{
		"flag":      cf.Name,
		"meta_name": f.qualifiedName(),
		"type":      f.typeString(),
		"required":  cf.Required,
		"alias":     cf.Aliases,
		"usage":     cf.Usage,
		"default":   cf.Value,
	}).Debug("append flag")
	efv.flags = append(efv.flags, cf)
}

func (efv *extractFlagsVisitor) visitBoolean(f *configField) {
	flag, _ := getFieldFlag(f)
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

	log.WithFields(log.Fields{
		"flag":      cf.Name,
		"meta_name": f.qualifiedName(),
		"type":      f.typeString(),
		"required":  cf.Required,
		"alias":     cf.Aliases,
		"usage":     cf.Usage,
		"default":   cf.Value,
	}).Debug("append flag")
	efv.flags = append(efv.flags, cf)
}

func (efv *extractFlagsVisitor) visitInt(f *configField) {
	flag, _ := getFieldFlag(f)
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

	log.WithFields(log.Fields{
		"flag":      cf.Name,
		"meta_name": f.qualifiedName(),
		"type":      f.typeString(),
		"required":  cf.Required,
		"alias":     cf.Aliases,
		"usage":     cf.Usage,
		"default":   cf.Value,
	}).Debug("append flag")
	efv.flags = append(efv.flags, cf)
}

type setConfigVisitor struct {
	ctx  *cli.Context
	errs []error
}

func newSetConfigVisitor(ctx *cli.Context) *setConfigVisitor {
	return &setConfigVisitor{ctx, nil}
}

var _ configVisitor = (*setConfigVisitor)(nil)

func (scv *setConfigVisitor) visit(f *configField) {
	if !f.isExported() {
		log.Debugf("config field '%s' is not exported", f.qualifiedName())
		return
	}

	if !f.hasTag() {
		log.Debugf("config field '%s' has no tag", f.qualifiedName())
		return
	}

	if _, err := getFieldFlag(f); err != nil {
		if _, err := getFieldArgPosition(f); err != nil {
			log.WithField("error", err).Debugf("config field '%s' has no flag or positional argument", f.qualifiedName())
		}
	}

	var err error
	switch f.fieldKind() {
	case reflect.String:
		err = scv.setString(f)
	case reflect.Bool:
		err = scv.setBoolean(f)
	case reflect.Int:
		err = scv.setInt(f)
	case reflect.Slice:
		err = scv.setSlice(f)
	default:
		panic("unsupported field kind")
	}
	if err != nil {
		scv.errs = append(scv.errs, err)
	}
}

func (scv *setConfigVisitor) setSlice(f *configField) error {
	et := f.fieldType().Elem()
	var err error
	switch et.Kind() {
	case reflect.String:
		err = scv.setStringSlice(f)
	default:
		panic(fmt.Errorf("unsupported slice kind '%s'", et.Kind()))
	}
	return err
}

func (scv *setConfigVisitor) setStringSlice(f *configField) error {
	var value []string
	ctx := scv.ctx
	if flag, err := getFieldFlag(f); err == nil {
		value = ctx.StringSlice(flag)
	}
	log.WithFields(log.Fields{"field": f.qualifiedName(), "value": value}).Debug("set []string")
	f.value.Set(reflect.ValueOf(value))
	return nil
}

func (scv *setConfigVisitor) setString(f *configField) error {
	var value string
	ctx := scv.ctx
	if flag, err := getFieldFlag(f); err == nil {
		value = ctx.String(flag)
	} else {
		pos, _ := getFieldArgPosition(f)
		value = ctx.Args().Get(pos)
	}
	log.WithFields(log.Fields{"field": f.qualifiedName(), "value": value}).Debug("set string")
	f.value.SetString(value)
	return nil
}

func (scv *setConfigVisitor) setBoolean(f *configField) error {
	flag, _ := getFieldFlag(f)
	value := scv.ctx.Bool(flag)
	f.value.SetBool(value)
	log.WithFields(log.Fields{"field": f.qualifiedName(), "value": value}).Debug("set boolean")
	return nil
}

func (scv *setConfigVisitor) setInt(f *configField) error {
	var value int
	ctx := scv.ctx
	if flag, err := getFieldFlag(f); err == nil {
		value = ctx.Int(flag)
	} else {
		pos, _ := getFieldArgPosition(f)
		av := ctx.Args().Get(pos)
		iv, err := strconv.Atoi(av)
		if err != nil {
			return errors.Wrapf(err, "failed to set field, expected int type but got '%s'", av)
		}
		value = iv
	}
	log.WithFields(log.Fields{"field": f.qualifiedName(), "value": value}).Debug("set int")
	f.value.SetInt(int64(value))
	return nil
}

func (scv *setConfigVisitor) getErrors() []error {
	return scv.errs
}

// isFieldRequired check if field is required
func isFieldRequired(f *configField) bool {
	tag := f.structField.Tag
	if tag == "" {
		return false
	}

	vt := tag.Get("validate")
	if vt == "" {
		return false
	}

	return strings.Contains(vt, "required")
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
	tag := f.structField.Tag
	if flag := tag.Get("flag"); flag != "" {
		return flag, nil
	}
	return "", errors.New("empty flag")
}

func getFieldAlias(f *configField) ([]string, error) {
	tag := f.structField.Tag
	if alias := tag.Get("alias"); alias != "" {
		return strings.Split(alias, ","), nil
	}
	return nil, errors.New("empty alias")
}

func getFieldArgPosition(f *configField) (int, error) {
	tag := f.structField.Tag
	if arg := tag.Get("arg"); arg != "" {
		return strconv.Atoi(arg)
	}
	return 0, errors.New("empty arg")
}

func getFieldUsage(f *configField) (string, error) {
	tag := f.structField.Tag
	usage := tag.Get("usage")
	if usage == "" {
		return "", errors.New("empty help")
	}

	choices := getFieldChoices(f)
	if len(choices) == 0 {
		return usage, nil
	}
	return fmt.Sprintf("%s (options: %s)", usage, strings.Join(choices, ", ")), nil
}

func reflectValueToStringSlice(f reflect.Value) []string {
	l := f.Len()
	s := make([]string, l)
	for i := l - 1; i >= 0; i-- {
		s[i] = f.Index(i).String()
	}
	return s
}

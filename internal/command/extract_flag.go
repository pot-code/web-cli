package command

import (
	"fmt"
	"reflect"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type extractFlagsVisitor struct {
	flags []cli.Flag
}

func newExtractFlagsVisitor() *extractFlagsVisitor {
	return &extractFlagsVisitor{flags: make([]cli.Flag, 0)}
}

var _ visitor = &extractFlagsVisitor{}

func (efv *extractFlagsVisitor) accept(f *configField) {
	if !f.isExported() {
		log.Debugf("config field '%s' is not exported", f.name())
		return
	}

	if !f.hasTag() {
		log.Debugf("config field '%s' has no tag", f.name())
		return
	}

	if _, err := getFlag(f); err != nil {
		log.Debugf("config field '%s' has no flag", f.name())
		return
	}

	switch f.kind() {
	case reflect.String:
		efv.visitString(f)
	case reflect.Int:
		efv.visitInt(f)
	case reflect.Bool:
		efv.visitBoolean(f)
	case reflect.Slice:
		efv.visitSlice(f)
	default:
		panic(fmt.Errorf("unsupported kind '%s'", f.kind()))
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
	flag, _ := getFlag(f)
	cf := &cli.StringSliceFlag{
		Name: flag,
	}

	if f.hasDefaultValue() {
		cf.Value = cli.NewStringSlice(getReflectedStringSlice(f.value)...)
	}

	if isRequired(f) {
		cf.Required = true
	}

	if u, err := getUsage(f); err == nil {
		cf.Usage = u
	}

	if alias, err := getAlias(f); err == nil {
		cf.Aliases = alias
	}

	log.WithFields(log.Fields{
		"flag":      cf.Name,
		"meta_name": f.name(),
		"type":      f.typeString(),
		"required":  cf.Required,
		"alias":     cf.Aliases,
		"usage":     cf.Usage,
		"default":   cf.Value,
	}).Debug("append flag")
	efv.flags = append(efv.flags, cf)

}

func (efv *extractFlagsVisitor) visitString(f *configField) {
	flag, _ := getFlag(f)
	cf := &cli.StringFlag{
		Name: flag,
	}

	if f.hasDefaultValue() {
		cf.Value = f.value.String()
	}

	if isRequired(f) {
		cf.Required = true
	}

	if u, err := getUsage(f); err == nil {
		cf.Usage = u
	}

	if alias, err := getAlias(f); err == nil {
		cf.Aliases = alias
	}

	log.WithFields(log.Fields{
		"flag":      cf.Name,
		"meta_name": f.name(),
		"type":      f.typeString(),
		"required":  cf.Required,
		"alias":     cf.Aliases,
		"usage":     cf.Usage,
		"default":   cf.Value,
	}).Debug("append flag")
	efv.flags = append(efv.flags, cf)
}

func (efv *extractFlagsVisitor) visitBoolean(f *configField) {
	flag, _ := getFlag(f)
	cf := &cli.BoolFlag{
		Name: flag,
	}

	if f.hasDefaultValue() {
		cf.Value = f.value.Bool()
	}

	if isRequired(f) {
		cf.Required = true
	}

	if u, err := getUsage(f); err == nil {
		cf.Usage = u
	}

	if alias, err := getAlias(f); err == nil {
		cf.Aliases = alias
	}

	log.WithFields(log.Fields{
		"flag":      cf.Name,
		"meta_name": f.name(),
		"type":      f.typeString(),
		"required":  cf.Required,
		"alias":     cf.Aliases,
		"usage":     cf.Usage,
		"default":   cf.Value,
	}).Debug("append flag")
	efv.flags = append(efv.flags, cf)
}

func (efv *extractFlagsVisitor) visitInt(f *configField) {
	flag, _ := getFlag(f)
	cf := &cli.IntFlag{
		Name: flag,
	}

	if f.hasDefaultValue() {
		cf.Value = int(f.value.Int())
	}

	if isRequired(f) {
		cf.Required = true
	}

	if u, err := getUsage(f); err == nil {
		cf.Usage = u
	}

	if alias, err := getAlias(f); err == nil {
		cf.Aliases = alias
	}

	log.WithFields(log.Fields{
		"flag":      cf.Name,
		"meta_name": f.name(),
		"type":      f.typeString(),
		"required":  cf.Required,
		"alias":     cf.Aliases,
		"usage":     cf.Usage,
		"default":   cf.Value,
	}).Debug("append flag")
	efv.flags = append(efv.flags, cf)
}

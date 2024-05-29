package command

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/urfave/cli/v2"
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
	et := f.getFieldType().Elem()
	switch et.Kind() {
	case reflect.String:
		efv.visitStringSlice(f)
	default:
		panic(fmt.Errorf("unsupported slice kind '%s'", et.Kind()))
	}
	return nil
}

func (efv *extractFlagsVisitor) visitStringSlice(f *configField) error {
	flag := f.getFlag()
	if flag == "" {
		return nil
	}

	cf := &cli.StringSliceFlag{
		Name: flag,
	}

	if f.hasDefaultValue() {
		cf.Value = cli.NewStringSlice(reflectValueToStringSlice(f.value)...)
	}

	cf.Required = f.isRequired()
	cf.Usage = f.getUsage()
	cf.Aliases = f.getAlias()

	efv.flags = append(efv.flags, cf)
	return nil
}

func (efv *extractFlagsVisitor) visitString(f *configField) error {
	flag := f.getFlag()
	if flag == "" {
		return nil
	}

	cf := &cli.StringFlag{
		Name: flag,
	}

	if f.hasDefaultValue() {
		cf.Value = f.value.String()
	}

	cf.Required = f.isRequired()
	cf.Usage = f.getUsage()
	cf.Aliases = f.getAlias()

	efv.flags = append(efv.flags, cf)
	return nil
}

func (efv *extractFlagsVisitor) visitBoolean(f *configField) error {
	flag := f.getFlag()
	if flag == "" {
		return nil
	}

	cf := &cli.BoolFlag{
		Name: flag,
	}

	if f.hasDefaultValue() {
		cf.Value = f.value.Bool()
	}

	cf.Required = f.isRequired()
	cf.Usage = f.getUsage()
	cf.Aliases = f.getAlias()

	efv.flags = append(efv.flags, cf)
	return nil
}

func (efv *extractFlagsVisitor) visitInt(f *configField) error {
	flag := f.getFlag()
	if flag == "" {
		return nil
	}

	cf := &cli.IntFlag{
		Name: flag,
	}

	if f.hasDefaultValue() {
		cf.Value = int(f.value.Int())
	}

	cf.Required = f.isRequired()
	cf.Usage = f.getUsage()
	cf.Aliases = f.getAlias()

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
	et := f.getFieldType().Elem()
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
	if flag := f.getFlag(); flag != "" {
		value = ctx.StringSlice(flag)
	}
	f.value.Set(reflect.ValueOf(value))
	return nil
}

func (scv *setConfigVisitor) visitString(f *configField) error {
	var value string
	ctx := scv.ctx
	if flag := f.getFlag(); flag != "" {
		value = ctx.String(flag)
	} else {
		pos, _ := f.getArgPosition()
		value = ctx.Args().Get(pos)
	}
	f.value.SetString(value)
	return nil
}

func (scv *setConfigVisitor) visitBoolean(f *configField) error {
	flag := f.getFlag()
	value := scv.ctx.Bool(flag)
	f.value.SetBool(value)
	return nil
}

func (scv *setConfigVisitor) visitInt(f *configField) error {
	var value int
	ctx := scv.ctx
	if flag := f.getFlag(); flag != "" {
		value = ctx.Int(flag)
	} else {
		pos, _ := f.getArgPosition()
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

func reflectValueToStringSlice(f reflect.Value) []string {
	l := f.Len()
	s := make([]string, l)
	for i := l - 1; i >= 0; i-- {
		s[i] = f.Index(i).String()
	}
	return s
}

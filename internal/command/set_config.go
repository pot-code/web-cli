package command

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type setConfigVisitor struct {
	ctx  *cli.Context
	errs []error
}

func newSetConfigVisitor(ctx *cli.Context) *setConfigVisitor {
	return &setConfigVisitor{ctx, nil}
}

var _ visitor = &setConfigVisitor{}

func (scv *setConfigVisitor) accept(f *configField) {
	if !f.isExported() {
		log.Debugf("config field '%s' is not exported", f.name())
		return
	}

	if !f.hasTag() {
		log.Debugf("config field '%s' has no tag", f.name())
		return
	}

	if _, err := getFlag(f); err != nil {
		if _, err := getArgPosition(f); err != nil {
			log.WithField("error", err).Debugf("config field '%s' has no flag or positional argument", f.name())
		}
	}

	var err error
	switch f.kind() {
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
	if flag, err := getFlag(f); err == nil {
		value = ctx.StringSlice(flag)
	}
	log.WithFields(log.Fields{"field": f.name(), "value": value}).Debug("set []string")
	f.value.Set(reflect.ValueOf(value))
	return nil
}

func (scv *setConfigVisitor) setString(f *configField) error {
	var value string
	ctx := scv.ctx
	if flag, err := getFlag(f); err == nil {
		value = ctx.String(flag)
	} else {
		pos, _ := getArgPosition(f)
		value = ctx.Args().Get(pos)
	}
	log.WithFields(log.Fields{"field": f.name(), "value": value}).Debug("set string")
	f.value.SetString(value)
	return nil
}

func (scv *setConfigVisitor) setBoolean(f *configField) error {
	flag, _ := getFlag(f)
	value := scv.ctx.Bool(flag)
	f.value.SetBool(value)
	log.WithFields(log.Fields{"field": f.name(), "value": value}).Debug("set boolean")
	return nil
}

func (scv *setConfigVisitor) setInt(f *configField) error {
	var value int
	ctx := scv.ctx
	if flag, err := getFlag(f); err == nil {
		value = ctx.Int(flag)
	} else {
		pos, _ := getArgPosition(f)
		av := ctx.Args().Get(pos)
		iv, err := strconv.Atoi(av)
		if err != nil {
			return errors.Wrapf(err, "failed to set field, expected int type but got '%s'", av)
		}
		value = iv
	}
	log.WithFields(log.Fields{"field": f.name(), "value": value}).Debug("set int")
	f.value.SetInt(int64(value))
	return nil
}

func (scv *setConfigVisitor) getErrors() []error {
	return scv.errs
}

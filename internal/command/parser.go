package command

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
)

type flagField struct {
	kind      reflect.Kind
	fieldName string
	flagName  string
	alias     string
}

type flagParser struct {
	flags  []cli.Flag
	fields []*flagField
}

func newFlagParser() *flagParser {
	return &flagParser{}
}

func (p *flagParser) parse(config interface{}) error {
	if err := validateConfig(config); err != nil {
		return err
	}

	rv := reflect.ValueOf(config)
	rt := rv.Type().Elem()
	for _, f := range reflect.VisibleFields(rt) {
		p.parseField(&f)
	}
	return nil
}

func (p *flagParser) parseField(field *reflect.StructField) {
	if _, ok := field.Tag.Lookup("flag"); !ok {
		return
	}

	fk := field.Type.Kind()
	switch field.Type.Kind() {
	case reflect.String:
		p.parseString(field)
	case reflect.Int:
		p.parseInt(field)
	case reflect.Bool:
		p.parseBool(field)
	case reflect.Float64:
		p.parseFloat(field)
	case reflect.Slice:
		p.parseSlice(field)
	default:
		panic(fmt.Errorf("unsupported config field kind: %s", fk))
	}
}

func (p *flagParser) parseString(field *reflect.StructField) {
	p.flags = append(p.flags, &cli.StringFlag{
		Name:    field.Tag.Get("flag"),
		Usage:   field.Tag.Get("usage"),
		Aliases: []string{field.Tag.Get("alias")},
	})

	p.fields = append(p.fields, &flagField{
		kind:      reflect.String,
		fieldName: field.Name,
		flagName:  field.Tag.Get("flag"),
	})
}

func (p *flagParser) parseInt(field *reflect.StructField) {
	p.flags = append(p.flags, &cli.IntFlag{
		Name:    field.Tag.Get("flag"),
		Usage:   field.Tag.Get("usage"),
		Aliases: []string{field.Tag.Get("alias")},
	})

	p.fields = append(p.fields, &flagField{
		kind:      reflect.Int,
		fieldName: field.Name,
		flagName:  field.Tag.Get("flag"),
	})
}

func (p *flagParser) parseBool(field *reflect.StructField) {
	p.flags = append(p.flags, &cli.BoolFlag{
		Name:    field.Tag.Get("flag"),
		Usage:   field.Tag.Get("usage"),
		Aliases: []string{field.Tag.Get("alias")},
	})

	p.fields = append(p.fields, &flagField{
		kind:      reflect.Bool,
		fieldName: field.Name,
		flagName:  field.Tag.Get("flag"),
	})
}

func (p *flagParser) parseFloat(field *reflect.StructField) {
	p.flags = append(p.flags, &cli.Float64Flag{
		Name:    field.Tag.Get("flag"),
		Usage:   field.Tag.Get("usage"),
		Aliases: []string{field.Tag.Get("alias")},
	})

	p.fields = append(p.fields, &flagField{
		kind:      reflect.Float64,
		fieldName: field.Name,
		flagName:  field.Tag.Get("flag"),
	})
}

func (p *flagParser) parseSlice(field *reflect.StructField) {
	ek := field.Type.Elem().Kind()
	switch ek {
	case reflect.String:
		p.flags = append(p.flags, &cli.StringSliceFlag{
			Name:    field.Tag.Get("flag"),
			Usage:   field.Tag.Get("usage"),
			Aliases: []string{field.Tag.Get("alias")},
		})
		p.fields = append(p.fields, &flagField{
			kind:      reflect.Slice,
			fieldName: field.Name,
			flagName:  field.Tag.Get("flag"),
		})
	case reflect.Int:
		p.flags = append(p.flags, &cli.IntSliceFlag{
			Name:    field.Tag.Get("flag"),
			Usage:   field.Tag.Get("usage"),
			Aliases: []string{field.Tag.Get("alias")},
		})
		p.fields = append(p.fields, &flagField{
			kind:      reflect.Slice,
			fieldName: field.Name,
			flagName:  field.Tag.Get("flag"),
		})
	case reflect.Float64:
		p.flags = append(p.flags, &cli.Float64SliceFlag{
			Name:    field.Tag.Get("flag"),
			Usage:   field.Tag.Get("usage"),
			Aliases: []string{field.Tag.Get("alias")},
		})
		p.fields = append(p.fields, &flagField{
			kind:      reflect.Slice,
			fieldName: field.Name,
			flagName:  field.Tag.Get("flag"),
		})
	default:
		panic(fmt.Errorf("unsupported config field kind: %s", ek))
	}
}

func (p *flagParser) getFlags() []cli.Flag {
	return p.flags
}

func (p *flagParser) getFields() []*flagField {
	return p.fields
}

type argField struct {
	kind      reflect.Kind
	position  int
	fieldName string
	alias     string
}

type argParser struct {
	fields []*argField
}

func newArgParser() *argParser {
	return &argParser{}
}

func (a *argParser) parse(config interface{}) error {
	if err := validateConfig(config); err != nil {
		return err
	}

	rv := reflect.ValueOf(config)
	rt := rv.Type().Elem()
	for _, f := range reflect.VisibleFields(rt) {
		a.parseField(&f)
	}
	return nil
}

func (a *argParser) parseField(field *reflect.StructField) {
	if _, ok := field.Tag.Lookup("arg"); !ok {
		return
	}

	fk := field.Type.Kind()
	switch field.Type.Kind() {
	case reflect.String:
		a.parseString(field)
	case reflect.Int:
		a.parseInt(field)
	case reflect.Bool:
		a.parseBool(field)
	case reflect.Float64:
		a.parseFloat64(field)
	default:
		panic(fmt.Errorf("unsupported config field kind: %s", fk))
	}
}

func (p *argParser) parseString(field *reflect.StructField) {
	pos := p.parsePosition(field)
	p.fields = append(p.fields, &argField{
		kind:      reflect.String,
		fieldName: field.Name,
		position:  pos,
		alias:     field.Tag.Get("alias"),
	})
}

func (p *argParser) parseInt(field *reflect.StructField) {
	pos := p.parsePosition(field)
	p.fields = append(p.fields, &argField{
		kind:      reflect.Int,
		fieldName: field.Name,
		position:  pos,
		alias:     field.Tag.Get("alias"),
	})
}

func (p *argParser) parseBool(field *reflect.StructField) {
	pos := p.parsePosition(field)
	p.fields = append(p.fields, &argField{
		kind:      reflect.Bool,
		fieldName: field.Name,
		position:  pos,
		alias:     field.Tag.Get("alias"),
	})
}

func (p *argParser) parseFloat64(field *reflect.StructField) {
	pos := p.parsePosition(field)
	p.fields = append(p.fields, &argField{
		kind:      reflect.Float64,
		fieldName: field.Name,
		position:  pos,
		alias:     field.Tag.Get("alias"),
	})
}

func (p *argParser) parsePosition(field *reflect.StructField) int {
	tv := field.Tag.Get("arg")
	pos, err := strconv.Atoi(tv)
	if err != nil {
		panic(fmt.Errorf("arg tag value must be int: %s", tv))
	}
	return pos
}

func (a *argParser) getFields() []*argField {
	return a.fields
}

func (a *argParser) getArgsUsage() string {
	var sb strings.Builder
	for _, f := range a.fields {
		sb.WriteString(fmt.Sprintf("%s ", f.alias))
	}
	return sb.String()
}

type argValueParser struct{}

func newArgValueParser() *argValueParser {
	return &argValueParser{}
}

func (a *argValueParser) parse(af *argField, v string) (any, error) {
	switch af.kind {
	case reflect.String:
		return v, nil
	case reflect.Int:
		return a.parseInt(v)
	case reflect.Bool:
		return a.parseBool(v)
	case reflect.Float64:
		return a.parseFloat64(v)
	default:
	}
	return nil, nil
}

func (a *argValueParser) parseInt(s string) (int, error) {
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("parse int: %w", err)
	}
	return v, nil
}

func (a *argValueParser) parseFloat64(s string) (float64, error) {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("parse float64: %w", err)
	}
	return v, nil
}

func (a *argValueParser) parseBool(s string) (bool, error) {
	v, err := strconv.ParseBool(s)
	if err != nil {
		return false, fmt.Errorf("parse bool: %w", err)
	}
	return v, nil
}

package command

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
)

type flagField struct {
	fieldKind reflect.Kind
	fieldName string
	flagName  string
}

type flagParser struct {
	flags  []cli.Flag
	fields []*flagField
}

func newFlagParser() *flagParser {
	return &flagParser{}
}

func (p *flagParser) parse(config any) error {
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
		panic(fmt.Errorf("不支持的配置类型: %s", fk))
	}
}

func (p *flagParser) parseString(field *reflect.StructField) {
	p.flags = append(p.flags, &cli.StringFlag{
		Name:    field.Tag.Get("flag"),
		Usage:   field.Tag.Get("usage"),
		Aliases: []string{field.Tag.Get("alias")},
	})

	p.fields = append(p.fields, &flagField{
		fieldKind: reflect.String,
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
		fieldKind: reflect.Int,
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
		fieldKind: reflect.Bool,
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
		fieldKind: reflect.Float64,
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
			fieldKind: reflect.Slice,
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
			fieldKind: reflect.Slice,
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
			fieldKind: reflect.Slice,
			fieldName: field.Name,
			flagName:  field.Tag.Get("flag"),
		})
	default:
		panic(fmt.Errorf("不支持的配置类型: %s", ek))
	}
}

type argField struct {
	fieldKind reflect.Kind
	fieldName string
	position  int
	alias     string // positional 参数的别名，一般用在命令行的 usage 里，全大写格式
}

type argParser struct {
	fields []*argField
}

func newArgParser() *argParser {
	return &argParser{}
}

func (a *argParser) parse(config any) error {
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
		panic(fmt.Errorf("不支持的配置类型: %s", fk))
	}
}

func (p *argParser) parseString(field *reflect.StructField) {
	pos := p.parsePosition(field)
	p.fields = append(p.fields, &argField{
		fieldKind: reflect.String,
		fieldName: field.Name,
		position:  pos,
		alias:     field.Tag.Get("alias"),
	})
}

func (p *argParser) parseInt(field *reflect.StructField) {
	pos := p.parsePosition(field)
	p.fields = append(p.fields, &argField{
		fieldKind: reflect.Int,
		fieldName: field.Name,
		position:  pos,
		alias:     field.Tag.Get("alias"),
	})
}

func (p *argParser) parseBool(field *reflect.StructField) {
	pos := p.parsePosition(field)
	p.fields = append(p.fields, &argField{
		fieldKind: reflect.Bool,
		fieldName: field.Name,
		position:  pos,
		alias:     field.Tag.Get("alias"),
	})
}

func (p *argParser) parseFloat64(field *reflect.StructField) {
	pos := p.parsePosition(field)
	p.fields = append(p.fields, &argField{
		fieldKind: reflect.Float64,
		fieldName: field.Name,
		position:  pos,
		alias:     field.Tag.Get("alias"),
	})
}

// 解析 positional arg 的位置
func (p *argParser) parsePosition(field *reflect.StructField) int {
	tv := field.Tag.Get("arg")
	pos, err := strconv.Atoi(tv)
	if err != nil {
		panic(fmt.Errorf("arg 标签值必须是整数: %s", tv))
	}
	if pos < 0 {
		panic(fmt.Errorf("arg 标签值不能小于0: %s", tv))
	}
	return pos
}

func (a *argParser) getArgsUsage() string {
	// 按照位置排序
	sort.Slice(a.fields, func(i, j int) bool {
		return a.fields[i].position < a.fields[j].position
	})

	var alias []string
	for _, f := range a.fields {
		alias = append(alias, f.alias)
	}
	return strings.Join(alias, " ")
}

type argValueParser struct{}

func newArgValueParser() *argValueParser {
	return &argValueParser{}
}

func (a *argValueParser) parse(af *argField, v string) (any, error) {
	switch af.fieldKind {
	case reflect.String:
		return v, nil
	case reflect.Int:
		return a.parseInt(v)
	case reflect.Bool:
		return a.parseBool(v)
	case reflect.Float64:
		return a.parseFloat64(v)
	default:
		panic(fmt.Errorf("不支持的配置类型: %s", af.fieldKind))
	}
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

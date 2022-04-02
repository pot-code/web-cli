package util

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/pot-code/web-cli/internal/validate"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type CommandFeature interface {
	Cond(c *cli.Context) bool
	Handle(c *cli.Context, cfg interface{}) error
}

type ConfigStructVisitor interface {
	VisitStringType(reflect.StructField, reflect.Value, *cli.Context) error
	VisitBooleanType(reflect.StructField, reflect.Value, *cli.Context) error
	VisitIntType(reflect.StructField, reflect.Value, *cli.Context) error
}

type CliCommand struct {
	cmd      *cli.Command
	cfg      interface{}
	services []CommandFeature
}

type CommandOption interface {
	apply(*cli.Command)
}

type optionFunc func(*cli.Command)

func (o optionFunc) apply(c *cli.Command) {
	o(c)
}

func WithAlias(alias []string) CommandOption {
	return optionFunc(func(c *cli.Command) {
		c.Aliases = alias
	})
}

func WithArgUsage(usage string) CommandOption {
	return optionFunc(func(c *cli.Command) {
		c.ArgsUsage = usage
	})
}

func WithFlags(flags cli.Flag) CommandOption {
	return optionFunc(func(c *cli.Command) {
		c.Flags = append(c.Flags, flags)
	})
}

func NewCliCommand(name, usage string, defaultConfig interface{}, options ...CommandOption) *CliCommand {
	cmd := &cli.Command{
		Name:  name,
		Usage: usage,
	}

	visitor := NewExtractFlagsVisitor()
	IterateCliConfig(defaultConfig, visitor, nil)
	cmd.Flags = append(cmd.Flags, visitor.Flags...)

	for _, option := range options {
		option.apply(cmd)
	}
	return &CliCommand{cmd, defaultConfig, []CommandFeature{}}
}

func (cc *CliCommand) AddFeature(feats ...CommandFeature) *CliCommand {
	cc.services = append(cc.services, feats...)
	return cc
}

func (cc *CliCommand) ExportCommand() *cli.Command {
	cc.cmd.Action = func(c *cli.Context) error {
		cfg := cc.cfg

		if cfg != nil {
			IterateCliConfig(cfg, NewSetCliValueVisitor(), c)
			if err := validate.V.Struct(cfg); err != nil {
				if v, ok := err.(validator.ValidationErrors); ok {
					msg := v[0].Translate(validate.T)
					cli.ShowCommandHelp(c, c.Command.Name)
					return NewCommandError(c.Command.Name, errors.New(msg))
				}
				return err
			}
		}

		for _, s := range cc.services {
			if s.Cond(c) {
				if err := s.Handle(c, cfg); err != nil {
					return err
				}
			}
		}

		return nil
	}
	return cc.cmd
}

func IterateCliConfig(config interface{}, visitor ConfigStructVisitor, runtime *cli.Context) {
	if config == nil {
		return
	}

	t := reflect.TypeOf(config)
	if t.Kind() != reflect.Ptr {
		panic("config must be of pointer type")
	}

	t = t.Elem()
	v := reflect.ValueOf(config)
	v = reflect.Indirect(v)
	for i := v.NumField() - 1; i >= 0; i-- {
		tf := t.Field(i)
		vf := v.Field(i)

		if !vf.CanSet() {
			log.WithFields(log.Fields{
				"caller": "ParseCliConfig",
			}).Warnf("config field [ %s ] can't be set, maybe it's not exported", tf.Name)
			continue
		}

		switch tf.Type.Kind() {
		case reflect.String:
			visitor.VisitStringType(tf, vf, runtime)
		case reflect.Bool:
			visitor.VisitBooleanType(tf, vf, runtime)
		case reflect.Int:
			visitor.VisitIntType(tf, vf, runtime)
		default:
			panic("not implemented")
		}
	}
}

type ExtractFlagsVisitor struct {
	Flags []cli.Flag
}

func NewExtractFlagsVisitor() *ExtractFlagsVisitor {
	return &ExtractFlagsVisitor{Flags: []cli.Flag{}}
}

var _ ConfigStructVisitor = &ExtractFlagsVisitor{}

func (efv *ExtractFlagsVisitor) VisitStringType(tf reflect.StructField, vf reflect.Value, c *cli.Context) error {
	tag := tf.Tag

	if flagName := tag.Get("flag"); flagName != "" {
		alias := tag.Get("alias")
		usage := tag.Get("usage")

		required := false
		if vf.IsZero() && validate.IsRequired(tag) {
			required = true
		}

		flag := &cli.StringFlag{
			Name:     flagName,
			Required: required,
		}

		if alias != "" {
			flag.Aliases = strings.Split(alias, ",")
		}

		if !vf.IsZero() {
			flag.Value = vf.String()
		}

		options := validate.GetOneOfItems(tag)
		if len(options) > 0 {
			usage += fmt.Sprintf(" (options: %s)", strings.Join(options, ", "))
		}
		flag.Usage = usage

		log.WithFields(log.Fields{
			"name":     flagName,
			"type":     tf.Type.String(),
			"required": required,
			"alias":    alias,
			"usage":    usage,
		}).Debug("register flag")

		efv.Flags = append(efv.Flags, flag)
	}
	return nil
}

func (efv *ExtractFlagsVisitor) VisitBooleanType(tf reflect.StructField, vf reflect.Value, c *cli.Context) error {
	if flagName := tf.Tag.Get("flag"); flagName != "" {
		usage := tf.Tag.Get("usage")
		alias := tf.Tag.Get("alias")

		flag := &cli.BoolFlag{
			Name:  flagName,
			Usage: usage,
		}
		if alias != "" {
			flag.Aliases = strings.Split(alias, ",")
		}
		if !vf.IsZero() {
			flag.Value = vf.Bool()
		}

		log.WithFields(log.Fields{
			"name":  flagName,
			"type":  tf.Type.String(),
			"alias": alias,
			"usage": usage,
		}).Debug("register flag")
		efv.Flags = append(efv.Flags, flag)
	}
	return nil
}

func (efv *ExtractFlagsVisitor) VisitIntType(tf reflect.StructField, vf reflect.Value, c *cli.Context) error {
	if flagTag := tf.Tag.Get("flag"); flagTag != "" {
		usage := tf.Tag.Get("usage")
		alias := tf.Tag.Get("alias")

		flagName := &cli.IntFlag{
			Name:  flagTag,
			Usage: usage,
		}
		if alias != "" {
			flagName.Aliases = strings.Split(alias, ",")
		}
		if !vf.IsZero() {
			flagName.Value = int(vf.Int())
		}

		log.WithFields(log.Fields{
			"name":  flagName,
			"type":  tf.Type.String(),
			"alias": alias,
			"usage": usage,
		}).Debug("register flag")
		efv.Flags = append(efv.Flags, flagName)
	}
	return nil
}

type SetCliValueVisitor struct{}

func NewSetCliValueVisitor() *SetCliValueVisitor {
	return &SetCliValueVisitor{}
}

var _ ConfigStructVisitor = &SetCliValueVisitor{}

func (efv *SetCliValueVisitor) VisitStringType(tf reflect.StructField, vf reflect.Value, c *cli.Context) error {
	if arg := tf.Tag.Get("arg"); arg != "" {
		index, err := strconv.Atoi(arg)
		if err != nil {
			panic(fmt.Sprintf("failed to convert [ %s ] to number", arg))
		}

		kind := tf.Type.Kind()
		if kind == reflect.String {
			vf.SetString(c.Args().Get(index))
		} else {
			panic(fmt.Errorf("unsupported arg type: %s", kind.String()))
		}
	} else if flag := tf.Tag.Get("flag"); flag != "" {
		vf.SetString(c.String(flag))
	}
	return nil
}

func (efv *SetCliValueVisitor) VisitBooleanType(tf reflect.StructField, vf reflect.Value, c *cli.Context) error {
	if flag := tf.Tag.Get("flag"); flag != "" {
		vf.SetBool(c.Bool(flag))
	}
	return nil
}

func (efv *SetCliValueVisitor) VisitIntType(tf reflect.StructField, vf reflect.Value, c *cli.Context) error {
	if flag := tf.Tag.Get("flag"); flag != "" {
		vf.SetInt(int64(c.Int("flag")))
	}
	return nil
}

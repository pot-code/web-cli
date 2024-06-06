package command

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/pot-code/web-cli/internal/validate"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

type CommandBuilder[T any] struct {
	name          string
	usage         string
	options       []CommandOption
	defaultConfig T
	handlers      []CommandHandler[T]
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

func NewBuilder[T any](name, usage string, defaultConfig T, options ...CommandOption) *CommandBuilder[T] {
	validateConfig(defaultConfig)

	return &CommandBuilder[T]{
		name:          name,
		usage:         usage,
		defaultConfig: defaultConfig,
		options:       options,
	}
}

func (cb *CommandBuilder[T]) AddHandlers(h ...CommandHandler[T]) *CommandBuilder[T] {
	cb.handlers = append(cb.handlers, h...)
	return cb
}

func (cb *CommandBuilder[T]) Build() *cli.Command {
	dc := cb.defaultConfig
	efv := newExtractFlagsVisitor()
	walkConfig(dc, efv)

	cmd := &cli.Command{
		Name:  cb.name,
		Usage: cb.usage,
	}

	for _, o := range cb.options {
		o.apply(cmd)
	}

	cmd.Flags = append(cmd.Flags, efv.getFlags()...)
	cmd.Before = func(c *cli.Context) error {
		scv := newSetConfigVisitor(c)
		if err := walkConfig(dc, scv); err != nil {
			return err
		}

		if err := validate.V.Struct(dc); err != nil {
			if v, ok := err.(validator.ValidationErrors); ok {
				msg := v[0].Translate(validate.T)
				cli.ShowCommandHelp(c, c.Command.Name)
				return errors.New(msg)
			}
			return err
		}
		return nil
	}
	cmd.Action = func(c *cli.Context) error {
		log.Debug().
			Str("config", fmt.Sprintf("%+v", dc)).
			Str("command", c.Command.FullName()).
			Msg("run command")
		for _, s := range cb.handlers {
			if err := s.Handle(c, dc); err != nil {
				return err
			}
		}
		return nil
	}
	return cmd
}

func validateConfig(config interface{}) {
	if config == nil {
		panic("config is nil")
	}

	if reflect.TypeOf(config).Kind() != reflect.Ptr {
		panic("config value must be of pointer type")
	}

	et := reflect.TypeOf(config).Elem()
	if et.Kind() != reflect.Struct {
		panic("config must be of struct type")
	}
}

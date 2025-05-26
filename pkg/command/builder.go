package command

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/pot-code/web-cli/pkg/validate"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

type CommandHandler[T any] interface {
	Handle(c *cli.Context, cfg T) error
}

type InlineHandler[T any] func(c *cli.Context, cfg T) error

func (s InlineHandler[T]) Handle(c *cli.Context, cfg T) error {
	return s(c, cfg)
}

type CommandBuilder[T any] struct {
	name          string
	usage         string
	defaultConfig T
	options       []commandOption
	handlers      []CommandHandler[T]
}

type commandOption func(*cli.Command)

func WithAlias(alias []string) commandOption {
	return func(c *cli.Command) {
		c.Aliases = alias
	}
}

func WithFlags(flags cli.Flag) commandOption {
	return func(c *cli.Command) {
		c.Flags = append(c.Flags, flags)
	}
}

func NewCommand[T any](name, usage string, defaultConfig T, options ...commandOption) *CommandBuilder[T] {
	return &CommandBuilder[T]{
		name:          name,
		usage:         usage,
		defaultConfig: defaultConfig,
		options:       options,
	}
}

func (cb *CommandBuilder[T]) AddHandler(h CommandHandler[T]) *CommandBuilder[T] {
	cb.handlers = append(cb.handlers, h)
	return cb
}

func (cb *CommandBuilder[T]) Create() *cli.Command {
	config := cb.defaultConfig
	if err := validateConfigValue(config); err != nil {
		panic(fmt.Errorf("validate config: %w", err))
	}

	fp := newFlagParser()
	if err := fp.parse(config); err != nil {
		panic(fmt.Errorf("parse flags: %w", err))
	}

	ap := newArgParser()
	if err := ap.parse(config); err != nil {
		panic(fmt.Errorf("parse args: %w", err))
	}

	cmd := &cli.Command{
		Name:  cb.name,
		Usage: cb.usage,
	}
	cmd.ArgsUsage = " " + ap.getArgsUsage()
	cmd.Flags = append(cmd.Flags, fp.flags...)

	for _, opt := range cb.options {
		opt(cmd)
	}

	cmd.Before = func(c *cli.Context) error {
		if !isConfigRequired(config) {
			return nil
		}

		// 从 cli 上下文中解析配置
		cp := newConfigParser(fp.fields, ap.fields)
		if err := cp.parseFromCliContext(c, config); err != nil {
			log.Err(err).Msg("set config from context")
			return err
		}

		// 校验配置
		if err := validate.V.Struct(config); err != nil {
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
			Str("config", fmt.Sprintf("%+v", config)).
			Str("command", c.Command.FullName()).
			Msg("run command")
		for _, s := range cb.handlers {
			if err := s.Handle(c, config); err != nil {
				return err
			}
		}
		return nil
	}

	return cmd
}

func isConfigRequired(config any) bool {
	rv := reflect.ValueOf(config)
	return rv.Kind() == reflect.Ptr && rv.IsNil()
}

func validateConfigValue(config any) error {
	if !isConfigRequired(config) {
		return nil
	}

	rv := reflect.ValueOf(config)
	if rv.Kind() != reflect.Ptr && rv.Elem().Kind() != reflect.Struct {
		return errors.New("config must be of pointer of struct type")
	}
	return nil
}

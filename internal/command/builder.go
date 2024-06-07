package command

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/pot-code/web-cli/internal/validate"
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
	options       []CommandOption
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

func WithFlags(flags cli.Flag) CommandOption {
	return optionFunc(func(c *cli.Command) {
		c.Flags = append(c.Flags, flags)
	})
}

func NewBuilder[T any](name, usage string, defaultConfig T, options ...CommandOption) *CommandBuilder[T] {
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
	config := cb.defaultConfig
	fp := newFlagParser()
	ap := newArgParser()

	if err := fp.parse(config); err != nil {
		panic(fmt.Errorf("parse flags: %w", err))
	}
	if err := ap.parse(config); err != nil {
		panic(fmt.Errorf("parse args: %w", err))
	}

	cmd := &cli.Command{
		Name:  cb.name,
		Usage: cb.usage,
	}
	for _, o := range cb.options {
		o.apply(cmd)
	}
	cmd.ArgsUsage = " " + ap.getArgsUsage()
	cmd.Flags = append(cmd.Flags, fp.getFlags()...)
	cmd.Before = func(c *cli.Context) error {
		cs := newConfigSetter(fp.getFields(), ap.getFields())
		if err := cs.setFromContext(c, config); err != nil {
			log.Err(err).Msg("set config from context")
			return err
		}

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

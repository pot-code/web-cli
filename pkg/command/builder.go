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

type CliCommand struct {
	cmd           *cli.Command
	defaultConfig interface{}
	handlers      []CommandHandler
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
	validateConfig(defaultConfig)
	cmd := &cli.Command{
		Name:  name,
		Usage: usage,
	}
	for _, option := range options {
		option.apply(cmd)
	}
	return &CliCommand{cmd, defaultConfig, []CommandHandler{}}
}

func (cc *CliCommand) AddHandlers(h ...CommandHandler) *CliCommand {
	cc.handlers = append(cc.handlers, h...)
	return cc
}

func (cc *CliCommand) BuildCommand() *cli.Command {
	dc := cc.defaultConfig
	efv := newExtractFlagsVisitor()
	walkConfig(dc, efv)
	cc.cmd.Flags = append(cc.cmd.Flags, efv.getFlags()...)

	cc.cmd.Before = func(c *cli.Context) error {
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

	cc.cmd.Action = func(c *cli.Context) error {
		log.Debug().
			Str("config", fmt.Sprintf("%+v", dc)).
			Str("command", c.Command.FullName()).
			Msg("run command")
		for _, s := range cc.handlers {
			if err := s.Handle(c, dc); err != nil {
				return err
			}
		}
		return nil
	}
	return cc.cmd
}

func validateConfig(config interface{}) {
	if config == nil {
		panic("config is nil")
	}

	if !isPointerType(config) {
		panic("config value must be of pointer type")
	}

	configType := reflect.TypeOf(config).Elem()
	if configType.Kind() != reflect.Struct {
		panic("config must be of struct type")
	}
}

func isPointerType(v interface{}) bool {
	if v == nil {
		return false
	}
	return reflect.TypeOf(v).Kind() == reflect.Ptr
}

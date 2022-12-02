package command

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/pot-code/web-cli/internal/validate"
	"github.com/urfave/cli/v2"
)

type CliCommand struct {
	cmd      *cli.Command
	cfg      interface{}
	handlers []CommandHandler
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
	cfg := cc.cfg

	efv := newExtractFlagsVisitor()
	walkConfig(cfg, efv)
	cc.cmd.Flags = append(cc.cmd.Flags, efv.getFlags()...)

	cc.cmd.Before = func(c *cli.Context) error {
		scv := newSetConfigVisitor(c)
		walkConfig(cfg, scv)

		if errs := scv.getErrors(); errs != nil {
			return errs[0]
		}

		if err := validate.V.Struct(cfg); err != nil {
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
		for _, s := range cc.handlers {
			if err := s.Handle(c, cfg); err != nil {
				return err
			}
		}
		return nil
	}
	return cc.cmd
}

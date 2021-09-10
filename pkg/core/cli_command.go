package core

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/pot-code/web-cli/pkg/validate"
	"github.com/urfave/cli/v2"
)

type CommandService interface {
	Cond(c *cli.Context) bool
	Handle(c *cli.Context, cfg interface{}) error
}

type ConfigStructVisitor interface {
	VisitStringType(reflect.StructField, reflect.Value, *cli.Context) error
	VisitBooleanType(reflect.StructField, reflect.Value, *cli.Context) error
	VisitIntType(reflect.StructField, reflect.Value, *cli.Context) error
}

type CliLeafCommand struct {
	cmd      *cli.Command
	cfg      interface{}
	services []CommandService
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

func NewCliLeafCommand(name, usage string, cfg interface{}, options ...CommandOption) *CliLeafCommand {
	cmd := &cli.Command{
		Name:  name,
		Usage: usage,
	}

	visitor := NewExtractFlagsVisitor()
	IterateCliConfig(cfg, visitor, nil)
	cmd.Flags = append(cmd.Flags, visitor.Flags...)

	for _, option := range options {
		option.apply(cmd)
	}
	return &CliLeafCommand{cmd, cfg, []CommandService{}}
}

func (cc *CliLeafCommand) AddService(services ...CommandService) *CliLeafCommand {
	cc.services = append(cc.services, services...)
	return cc
}

func (cc *CliLeafCommand) ExportCommand() *cli.Command {
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

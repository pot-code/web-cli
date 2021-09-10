package core

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/pot-code/web-cli/pkg/validate"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

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

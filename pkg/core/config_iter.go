package core

import (
	"reflect"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

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

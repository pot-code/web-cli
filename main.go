package main

import (
	"os"

	"github.com/pkg/errors"
	"github.com/pot-code/web-cli/cmd"
	"github.com/pot-code/web-cli/pkg/util"
	log "github.com/sirupsen/logrus"
)

//go:generate go install github.com/valyala/quicktemplate/qtc
//go:generate qtc -dir=templates
func main() {
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: true,
		QuoteEmptyFields: true,
	})
	if err := cmd.RootCmd.Run(os.Args); err != nil {
		if ste, ok := err.(util.StackTracer); ok {
			trace := util.GetVerboseStackTrace(3, ste)
			cause := errors.Cause(err)
			log.WithFields(log.Fields{"error": cause.Error(), "stack": trace}).Debug(err.Error())
		}
		log.Error(err.Error())
	}
}

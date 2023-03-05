package main

import (
	"embed"
	"os"

	"github.com/pot-code/web-cli/cmd"
	"github.com/pot-code/web-cli/pkg/provider"
	"github.com/pot-code/web-cli/pkg/validate"
	log "github.com/sirupsen/logrus"
)

//go:embed templates
var templates embed.FS

func main() {
	provider.InitTemplateFS(templates)
	validate.InitValidator()

	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: true,
		QuoteEmptyFields: true,
	})

	if err := cmd.RootCmd.Run(os.Args); err != nil {
		log.Error(err.Error())
	}
}

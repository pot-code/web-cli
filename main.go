package main

import (
	"embed"
	"os"

	"github.com/pot-code/web-cli/cmd"
	"github.com/pot-code/web-cli/internal/provider"
	"github.com/pot-code/web-cli/internal/validate"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//go:embed templates
var templates embed.FS

func main() {
	provider.InitTemplateFS(templates)
	validate.InitValidator()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if err := cmd.RootCmd.Run(os.Args); err != nil {
		log.Err(err).Msg("run command")
	}
}

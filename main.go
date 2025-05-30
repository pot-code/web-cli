package main

import (
	"embed"
	"os"

	"github.com/pot-code/web-cli/cmd"
	"github.com/pot-code/web-cli/pkg/provider"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//go:embed templates
var templates embed.FS

func main() {
	provider.InitTemplateFS(templates)

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if err := cmd.RootCmd.Run(os.Args); err != nil {
		log.Err(err).Msg("run command")
	}
}

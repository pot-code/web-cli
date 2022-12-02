package main

import (
	"os"

	"github.com/pot-code/web-cli/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: true,
		QuoteEmptyFields: true,
	})
	if err := cmd.RootCmd.Run(os.Args); err != nil {
		log.Error(err.Error())
	}
}

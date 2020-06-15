package main

import (
	"log"

	"github.com/pot-code/web-cli/cmd"
)

func main() {
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

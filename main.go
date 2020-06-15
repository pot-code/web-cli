package main

import (
	"log"

	"github.com/pot-code/web-cli/cmd"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

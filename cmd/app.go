package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "web-cli",
	Short: "CLI tool to generate FE/BE project based on template",
}

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get CWD: %s\n", err)
	}
	generate := NewGenerateCommand(cwd).Init()
	root.AddCommand(generate)
}

// Run execute root command
func Run() error {
	return root.Execute()
}

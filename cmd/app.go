package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "web-cli",
	Short: "CLI tool to generate FE/BE project based on template",
	Run: func(cmd *cobra.Command, args []string) {
		debug, _ := cmd.PersistentFlags().GetBool("debug")
		if debug {
			log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		}
	},
}

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get CWD: %s\n", err)
	}
	generate := NewGenerateCommand(cwd).Init()
	root.AddCommand(generate)
	root.PersistentFlags().Bool("verbose", false, "verbose output")
	root.PersistentFlags().Bool("debug", false, "debug output")
}

// Run execute root command
func Run() error {
	return root.Execute()
}

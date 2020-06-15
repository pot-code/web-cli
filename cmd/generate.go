package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"sync"
	"time"

	"github.com/pot-code/web-cli/template/backend"
	"github.com/spf13/cobra"
)

// prefix settings
const (
	FrontendPrefix = "frontend"
	BackendPrefix  = "backend"
	ConfigPrefix   = "config"
	MainPrefix     = "cmd"
)

// template URL
const (
	FrontendURL = "https://github.com/pot-code/react-boilerplate"
	BackendURL  = "https://github.com/pot-code/go-boilerplate"
)

// GenerateCommand TODO
type GenerateCommand struct {
	AppName   string
	Templates []*TemplateEntry
	cwd       string
	verbose   bool   // verbose outpout
	root      string // cwd/<app name>/
}

// NewGenerateCommand create GenerateCommand instance
func NewGenerateCommand(cwd string) *GenerateCommand {
	return &GenerateCommand{
		cwd: cwd,
	}
}

// Init init GenerateCommand
func (gc GenerateCommand) Init() *cobra.Command {
	// gc.registerEnvVariables()

	cmd := &cobra.Command{
		Use:   "generate NAME(project name)",
		Short: "generate an empty project based on templates",
		Args:  cobra.MinimumNArgs(1), // requires a name argument
		Run:   gc.run,
	}
	cmd.Flags().String("dirname", "", "project folder name")
	cmd.Flags().StringP("module", "M", "", "go module name (required)")
	cmd.Flags().String("env-prefix", "GO_WEB", "env variable prefix")
	cmd.Flags().String("config", "", "yaml config for additional config fields")
	cmd.Flags().StringP("desc-short", "D", "", "project binary short description")
	cmd.MarkFlagRequired("module")
	return cmd
}

// run run function for cobra command run
func (gc *GenerateCommand) run(cmd *cobra.Command, args []string) {
	debug, _ := cmd.Flags().GetBool("debug")
	if debug {
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	}

	verbose, _ := cmd.Flags().GetBool("verbose")
	gc.verbose = verbose

	gc.AppName = args[0]
	folderName, _ := cmd.Flags().GetString("dirname")
	if folderName != "" {
		gc.root = folderName
	} else {
		gc.root = path.Join(gc.cwd, gc.AppName)
	}

	log.Println("Loading template...")
	gc.registerTemplateEntry()

	// clone boilerplate
	ctx, cancel := context.WithCancel(context.Background())
	waitGroup := new(sync.WaitGroup)
	errch := make(chan error)
	waitDone := make(chan struct{})
	waitGroup.Add(2)
	go func() {
		// log.Println("Cloning FE repo...")
		if err := gc.cloneFrontend(ctx); err != nil {
			errch <- err
		}
		waitGroup.Done()
	}()
	go func() {
		// log.Println("Cloning BE repo...")
		if err := gc.cloneBackend(ctx); err != nil {
			errch <- err
		}
		waitGroup.Done()
	}()
	go func() {
		waitGroup.Wait()
		close(waitDone)
	}()

	select {
	case <-waitDone:
		cancel()
	case err := <-errch:
		cancel()
		log.Fatal(err)
	}

	log.Println("Generating template...")
	gc.generateTemplate(cmd, args)

	log.Println("Init module...")
	gc.initModule(cmd, args)

	log.Println("Clean up .git")
	gc.cleanGit(nil)
}

// registerTemplateEntry TODO
func (gc *GenerateCommand) registerTemplateEntry() {
	var entries []*TemplateEntry

	root := gc.root
	entries = append(entries,
		&TemplateEntry{
			Name:     "config",
			Template: backend.ConfigTemplate,
			Output:   path.Join(root, BackendPrefix, "config.go"),
		},
		&TemplateEntry{
			Name:     "main",
			Template: backend.MainTemplate,
			Output:   path.Join(root, BackendPrefix, "main.go"),
		},
	)
	gc.Templates = entries
}

// cloneFrontend clone frontend template from
// https://github.com/pot-code/react-boilerplate
func (gc GenerateCommand) cloneFrontend(ctx context.Context) error {
	// env := gc.Env
	destName := path.Join(gc.root, FrontendPrefix)
	if console, err := Clone(ctx, FrontendURL, []string{destName}); err != nil {
		log.Print(string(console))
		return fmt.Errorf("Error while cloning frontend repo: %w", err)
	} else if gc.verbose {
		log.Print(string(console))
	}
	return nil
}

// cloneBackend clone backend template from
// https://github.com/pot-code/go-boilerplate
func (gc GenerateCommand) cloneBackend(ctx context.Context) error {
	// env := gc.Env
	destName := path.Join(gc.root, BackendPrefix)
	if console, err := Clone(ctx, BackendURL, []string{destName}); err != nil {
		log.Print(string(console))
		return fmt.Errorf("Error while cloning backend repo: %w", err)
	} else if gc.verbose {
		log.Print(string(console))
	}
	return nil
}

// generateTemplate generate files from templates
func (gc GenerateCommand) generateTemplate(cmd *cobra.Command, args []string) {
	appName := gc.AppName
	moduleName, _ := cmd.Flags().GetString("module")
	envPrefix, _ := cmd.Flags().GetString("env-prefix")
	yamlConfig, _ := cmd.Flags().GetString("config")
	short, _ := cmd.Flags().GetString("desc-short")

	data := NewTemplateData(envPrefix, moduleName, appName)
	data.Short = short
	if yamlConfig != "" {
		log.Printf("Loading template data from '%s'...\n", yamlConfig)
		if err := data.LoadFromYaml(yamlConfig); err != nil {
			log.Fatal(err)
		}
	}

	entries := gc.Templates
	for _, entry := range entries {
		if err := CreateFromTemplate(entry, data); err != nil {
			log.Fatal(err)
		}
	}
}

// initModule run go mod init, go mod tidy, etc.
// TODO: extract go mod command to functions
func (gc GenerateCommand) initModule(cmd *cobra.Command, args []string) {
	moduleName, _ := cmd.Flags().GetString("module")
	cwd := path.Join(gc.root, BackendPrefix)
	verbose := gc.verbose

	// init
	init := exec.Command("go", "mod", "init", moduleName)
	init.Dir = cwd
	out, err := init.CombinedOutput()
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			log.Fatal("Go is not installed(https://golang.org/dl/) or not exists in PATH")
		}
		log.Fatalf("Error while doing 'go mod init': %s, %s\n", string(out), err)
	}
	if verbose {
		log.Print(string(out))
	}

	// format code(auto import)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if console, err := Goimports(ctx, cwd); err != nil {
		log.Print(string(console))
		log.Fatalf("Error while executing goimports: %s, %s\n", console, err)
	} else if verbose && len(console) > 0 {
		log.Print(string(console))
	}

	// tidy
	tidy := exec.Command("go", "mod", "tidy")
	tidy.Dir = cwd
	out, err = tidy.CombinedOutput()
	if err != nil {
		log.Fatalf("Error while doing 'go mod tidy': %s, %s\n", string(out), err)
	}
	if verbose {
		log.Print(string(out))
	}
}

func (gc GenerateCommand) cleanGit(ctx context.Context) error {
	frontend := path.Join(gc.root, FrontendPrefix, ".git")
	backend := path.Join(gc.root, BackendPrefix, ".git")

	if err := os.RemoveAll(frontend); err != nil {
		return fmt.Errorf("Error while removing frontend git dir: %w", err)
	}
	if err := os.RemoveAll(backend); err != nil {
		return fmt.Errorf("Error while removing backend git dir: %w", err)
	}
	return nil
}

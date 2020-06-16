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
	appName   string
	templates []*TemplateEntry
	repos     []*GithubRepository
	cwd       string
	verbose   bool   // verbose outpout
	root      string // cwd/<app name>/
}

// GithubRepository TODO
type GithubRepository struct {
	Name   string
	URL    string
	Output string // clone destination
}

// NewGenerateCommand create GenerateCommand instance
func NewGenerateCommand(cwd string) *GenerateCommand {
	return &GenerateCommand{
		cwd: cwd,
	}
}

// Init init GenerateCommand
func (gc GenerateCommand) Init() *cobra.Command {
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

	gc.appName = args[0]
	folderName, _ := cmd.Flags().GetString("dirname")
	if folderName != "" {
		gc.root = folderName
	} else {
		gc.root = path.Join(gc.cwd, gc.appName)
	}

	gc.registerTemplateEntry()
	gc.registerGithubRepo()

	// clone boilerplate
	log.Println("Cloning templates...")
	gc.cloneTemplates()

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
	gc.templates = entries
}

// registerGithubRepo TODO
func (gc *GenerateCommand) registerGithubRepo() {
	var repos []*GithubRepository

	root := gc.root
	repos = append(repos,
		&GithubRepository{
			Name:   "backend",
			URL:    BackendURL,
			Output: path.Join(root, BackendPrefix),
		},
		&GithubRepository{
			Name:   "frontend",
			URL:    FrontendURL,
			Output: path.Join(root, FrontendPrefix),
		},
	)
	gc.repos = repos
}

// cloneTemplates TODO
func (gc GenerateCommand) cloneTemplates() {
	ctx, cancel := context.WithCancel(context.Background())
	repos := gc.repos
	waitGroup := new(sync.WaitGroup)
	errch := make(chan error)
	waitDone := make(chan struct{})
	goroutines := make([]func() error, len(repos))

	waitGroup.Add(len(repos))
	for i, v := range repos {
		func(idx int, repo *GithubRepository) {
			goroutines[idx] = func() error {
				if console, err := Clone(ctx, repo.URL, []string{repo.Output}); err != nil {
					log.Print(string(console))
					return fmt.Errorf("Error while cloning %s repo: %w", repo.Name, err)
				} else if gc.verbose && console != nil {
					log.Print(string(console))
				}
				return nil
			}
		}(i, v)
	}
	for _, v := range goroutines {
		go func(fn func() error) {
			if err := fn(); err != nil {
				errch <- err
			}
			waitGroup.Done()
		}(v)
	}
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
}

// generateTemplate generate files from templates
func (gc GenerateCommand) generateTemplate(cmd *cobra.Command, args []string) {
	appName := gc.appName
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

	entries := gc.templates
	for _, entry := range entries {
		if err := CreateFromTemplate(entry, data); err != nil {
			log.Fatal(err)
		}
	}
}

// initModule run go mod init, go mod tidy, etc.
// TODO: extract go mod commands to functions
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
	} else if verbose && console != nil {
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

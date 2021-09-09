package cmd

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/pot-code/web-cli/pkg/constants"
	"github.com/pot-code/web-cli/pkg/core"
	"github.com/pot-code/web-cli/pkg/util"
	"github.com/pot-code/web-cli/templates"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"golang.org/x/tools/go/ast/astutil"
)

const (
	structHandlerCollection  = "HandlerCollection"
	varWireHandlerCollection = "WireHandlerCollection"
	varHttpSet               = "HttpSet"
	varWireSet               = "WireSet"
)

type genApiConfig struct {
	GenType         string `name:"type" validate:"required,oneof=go"`    // generation type
	ArgName         string `arg:"0" name:"NAME" validate:"required,var"` // go pkg name
	PackageName     string
	ProjectName     string
	CamelModuleName string
	AuthorName      string
}

var genAPICmd = &cli.Command{
	Name:      "api",
	Usage:     "generate an api module",
	ArgsUsage: "MODULE_NAME",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "type",
			Aliases: []string{"T"},
			Usage:   "api type (go)",
			Value:   "go",
		},
	},
	Action: func(c *cli.Context) error {
		cfg := new(genApiConfig)
		err := util.ParseConfig(c, cfg)
		if err != nil {
			if _, ok := err.(*util.CommandError); ok {
				cli.ShowCommandHelp(c, c.Command.Name)
			}
			return err
		}

		pkgName := strings.ReplaceAll(cfg.ArgName, "-", "_")
		pkgName = strcase.ToSnake(pkgName)
		cfg.PackageName = pkgName

		if cfg.GenType == "go" {
			meta, err := util.ParseGoMod("go.mod")
			if err != nil {
				return err
			}

			cfg.CamelModuleName = strcase.ToCamel(pkgName)
			log.Debug("generated module name: ", cfg.CamelModuleName)
			cfg.AuthorName = meta.Author
			cfg.ProjectName = meta.ProjectName

			dir := cfg.PackageName
			log.Debug("output path: ", dir)
			_, err = os.Stat(dir)
			if err == nil {
				log.Infof("[skipped]'%s' already exists", dir)
				return nil
			}

			registryFile := path.Join("server", "server.go")
			if err := updateServerRegistry(cfg, registryFile, newAddHandlerVisitor(cfg)); err != nil {
				return fmt.Errorf("failed to update registry: %w", err)
			}

			gen := generateGoApiFiles(cfg)
			if err := gen.Run(); err != nil {
				return err
			}
		}
		return nil
	},
}

func generateGoApiFiles(config *genApiConfig) core.Runner {
	modelName := config.CamelModuleName
	pkgName := config.PackageName
	handlerFile := fmt.Sprintf("%s_handler.go", pkgName)
	handlerName := fmt.Sprintf(constants.GoApiHandlerPattern, modelName)
	svcName := fmt.Sprintf(constants.GoApiServicePattern, modelName)
	repoName := fmt.Sprintf(constants.GoApiRepositoryPattern, modelName)

	return util.NewTaskComposer("").AddFile(
		&core.FileDesc{
			Path: path.Join("server", handlerFile),
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoApiHandler(&buf, config.ProjectName, config.AuthorName, pkgName, svcName, handlerName)
				return buf.Bytes()
			},
		},
		&core.FileDesc{
			Path: path.Join(pkgName, "model.go"),
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoApiModel(&buf, pkgName, svcName, repoName)
				return buf.Bytes()
			},
		},
		&core.FileDesc{
			Path: path.Join(pkgName, "repo.go"),
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoApiRepo(&buf, pkgName, repoName)
				return buf.Bytes()
			},
		},
		&core.FileDesc{
			Path: path.Join(pkgName, "service.go"),
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoApiService(&buf, pkgName, svcName, repoName)
				return buf.Bytes()
			},
		},
	).AddCommand(
		&core.Command{
			Bin:  "ent",
			Args: []string{"init", config.CamelModuleName},
		},
		&core.Command{
			Bin:  "wire",
			Args: []string{"./server"},
		},
		&core.Command{
			Bin:  "go",
			Args: []string{"mod", "tidy"},
		},
	)
}

func updateServerRegistry(cfg *genApiConfig, registryFile string, visitor serverRegistryVisitor) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, registryFile, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse '%s': %w", registryFile, err)
	}

	visitor.visitRoot(f, fset)
	modifedAst := astutil.Apply(f, func(c *astutil.Cursor) bool {
		n := c.Node()

		if _, ok := n.(*ast.File); ok {
			return true
		}

		if g, ok := n.(*ast.GenDecl); ok {
			return g.Tok == token.TYPE || g.Tok == token.VAR
		}

		if ts, ok := n.(*ast.TypeSpec); ok {
			return ts.Name.String() == structHandlerCollection
		}

		if vs, ok := n.(*ast.ValueSpec); ok {
			switch vs.Names[0].String() {
			case varHttpSet:
				visitor.visitHttpSet(vs)
			}
		}

		if st, ok := n.(*ast.StructType); ok {
			visitor.visitHandlerCollection(st)
		}

		return false
	}, nil)

	out, err := os.OpenFile(registryFile, os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to update '%s': %w", registryFile, err)
	}
	return format.Node(out, fset, modifedAst)
}

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
	"github.com/pkg/errors"
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

type GenApiConfig struct {
	GenType string `flag:"type" alias:"t" usage:"api type" validate:"required,oneof=go"` // generation type
	ArgName string `arg:"0" alias:"module_name" validate:"required,var"`                 // go pkg name
}

var GenAPICmd = core.NewCliLeafCommand("api", "generate an api module",
	&GenApiConfig{
		GenType: "go",
	},
	core.WithArgUsage("module_name"),
).AddService(
	&GenerateGoApiService{
		RegistryFile: path.Join("server", "server.go"),
	},
).ExportCommand()

type GenerateGoApiService struct {
	RegistryFile    string
	PackageName     string
	ProjectName     string
	CamelModuleName string
	AuthorName      string
	Config          *GenApiConfig
}

var _ core.CommandService = &GenerateGoApiService{}

func (gga *GenerateGoApiService) Cond(c *cli.Context) bool {
	return c.String("type") == "go"
}

func (gga *GenerateGoApiService) Handle(c *cli.Context, cfg interface{}) error {
	config := cfg.(*GenApiConfig)
	pkgName := strings.ReplaceAll(config.ArgName, "-", "_")
	pkgName = strcase.ToSnake(pkgName)
	gga.PackageName = pkgName

	meta, err := util.ParseGoMod(constants.GoModFile)
	if err != nil {
		return errors.WithStack(err)
	}

	gga.CamelModuleName = strcase.ToCamel(pkgName)
	log.Debug("generated module name: ", gga.CamelModuleName)
	gga.AuthorName = meta.Author
	gga.ProjectName = meta.ProjectName

	dir := gga.PackageName
	log.Debug("output path: ", dir)
	_, err = os.Stat(dir)
	if err == nil {
		log.Infof("[skipped]'%s' already exists", dir)
		return nil
	}

	if err := gga.updateServerRegistry(newAddHandlerVisitor(gga)); err != nil {
		return errors.Wrap(err, "failed to update registry")
	}

	return gga.generateFiles().Run()
}

func (gga *GenerateGoApiService) generateFiles() core.Runner {
	projectName := gga.ProjectName
	authorName := gga.AuthorName
	modelName := gga.CamelModuleName
	pkgName := gga.PackageName
	handlerFile := fmt.Sprintf("%s_handler.go", pkgName)
	handlerName := fmt.Sprintf(constants.GoApiHandlerPattern, modelName)
	svcName := fmt.Sprintf(constants.GoApiServicePattern, modelName)
	repoName := fmt.Sprintf(constants.GoApiRepositoryPattern, modelName)

	return util.NewTaskComposer("").AddFile(
		&core.FileDesc{
			Path: path.Join("server", handlerFile),
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteGoApiHandler(&buf, projectName, authorName, pkgName, svcName, handlerName)
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
		&core.ShellCommand{
			Bin:  "ent",
			Args: []string{"init", modelName},
		},
		&core.ShellCommand{
			Bin:  "wire",
			Args: []string{"./server"},
		},
		&core.ShellCommand{
			Bin:  "go",
			Args: []string{"mod", "tidy"},
		},
	)
}

func (gga *GenerateGoApiService) updateServerRegistry(visitor serverRegistryVisitor) error {
	registryFile := gga.RegistryFile
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

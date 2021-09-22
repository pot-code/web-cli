package cmd

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/pot-code/web-cli/pkg/constants"
	"golang.org/x/tools/go/ast/astutil"
)

type serverRegistryVisitor interface {
	visitHttpSet(node *ast.ValueSpec) error
	visitHandlerCollection(node *ast.StructType) error
	visitRoot(root *ast.File, fset *token.FileSet) error
}

type addHandlerVisitor struct {
	cfg         *GenerateGoApiService
	handlerName string
	// handlerConstructorName string
}

var _ serverRegistryVisitor = &addHandlerVisitor{}

func newAddHandlerVisitor(cfg *GenerateGoApiService) *addHandlerVisitor {
	handlerName := fmt.Sprintf(constants.GoApiHandlerPattern, cfg.CamelModuleName)
	// handlerConstructorName := constants.GoConstructorPrefix + handlerName
	// return &addHandlerVisitor{cfg, handlerName, handlerConstructorName}
	return &addHandlerVisitor{cfg, handlerName}
}

func (srv *addHandlerVisitor) visitHttpSet(node *ast.ValueSpec) error {
	// handlerConstructorName := srv.handlerConstructorName
	// pkgName := srv.cfg.PackageName
	// ce := node.Values[0].(*ast.CallExpr)

	// ce.Args = append(ce.Args,
	// 	&ast.Ident{Name: handlerConstructorName},
	// 	&ast.SelectorExpr{X: &ast.Ident{Name: pkgName}, Sel: &ast.Ident{Name: varWireSet}},
	// )
	return nil
}

func (srv *addHandlerVisitor) visitHandlerCollection(node *ast.StructType) error {
	pkgName := srv.cfg.PackageName
	handlerName := srv.handlerName

	node.Fields.List = append(node.Fields.List, &ast.Field{
		Names: []*ast.Ident{{Name: handlerName}},
		Type: &ast.StarExpr{X: &ast.SelectorExpr{
			X:   &ast.Ident{Name: pkgName},
			Sel: &ast.Ident{Name: handlerName},
		}},
	})
	return nil
}

func (srv *addHandlerVisitor) visitRoot(root *ast.File, fset *token.FileSet) error {
	cfg := srv.cfg
	pkgPath := fmt.Sprintf(`%s/%s/%s/%s`,
		constants.GoGithubPrefix,
		cfg.AuthorName,
		cfg.ProjectName,
		cfg.PackageName,
	)
	astutil.AddImport(fset, root, pkgPath)
	return nil
}

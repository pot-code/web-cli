package cmd

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"

	"github.com/fatih/structtag"
	"github.com/pot-code/web-cli/internal/validate"
)

type tagMeta struct {
	Default     string
	Required    bool
	MarshalName string
}

func parseTag(tag reflect.StructTag) (*tagMeta, error) {
	tt := strings.Trim(string(tag), "`")

	tags, err := structtag.Parse(tt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tag: %w", err)
	}

	tg := new(tagMeta)

	ms, err := tags.Get("mapstructure")
	if err == nil {
		tg.MarshalName = ms.Value()
	}

	d, err := tags.Get("default")
	if err == nil {
		tg.Default = d.Value()
	}

	tg.Required = validate.IsRequired(tag)
	return tg, nil
}

func parseConfigFile(config *GoFlagsConfig) (*configStructVisitor, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, config.ConfigPath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	pkg := f.Name.String()
	visitor := newConfigStructVisitor(fset, pkg)
	ast.Inspect(f, func(n ast.Node) bool {
		if _, ok := n.(*ast.File); ok {
			return true
		}

		if g, ok := n.(*ast.GenDecl); ok {
			return g.Tok == token.TYPE
		}

		if ts, ok := n.(*ast.TypeSpec); ok {
			return strings.HasSuffix(ts.Name.String(), "Config")
		}

		if st, ok := n.(*ast.StructType); ok {
			visitor.visit(st)
		}

		return false
	})

	return visitor, nil
}

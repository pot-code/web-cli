package cmd

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/fatih/structtag"
)

type tagMeta struct {
	Default     string
	Required    bool
	MarshalName string
}

func parseTag(tag string) (*tagMeta, error) {
	tt := strings.Trim(tag, "`")

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

	r, err := tags.Get("validate")
	if err == nil {
		options := strings.Split(r.Value(), ",")
		for _, o := range options {
			if o == "required" {
				tg.Required = true
				break
			}
		}
	}

	return tg, nil
}

func parseConfigFile(config *genFlagsConfig) (*configStructVisitor, error) {
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

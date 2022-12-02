package util

import (
	"go/ast"
	"go/parser"
	"go/token"

	"golang.org/x/tools/go/ast/astutil"
)

type GoAstModification interface {
	//Selector select target node
	Selector(c *astutil.Cursor) bool
	// Action do action on selected node
	Action(c *astutil.Cursor) error
}

type GoAstModifier struct {
	file    string
	fset    *token.FileSet
	actors  []GoAstModification
	imports []string
}

func NewGoAstModifier(file string) *GoAstModifier {
	return &GoAstModifier{file: file, actors: []GoAstModification{}, imports: make([]string, 0)}
}

func (gas *GoAstModifier) AddModification(actor ...GoAstModification) {
	gas.actors = append(gas.actors, actor...)
}

func (gas *GoAstModifier) AddImport(p string) {
	gas.imports = append(gas.imports, p)
}

func (gas *GoAstModifier) ParseAndModify() (fset *token.FileSet, n ast.Node, err error) {
	fset = token.NewFileSet()
	f, err := parser.ParseFile(fset, gas.file, nil, parser.ParseComments)
	if err != nil {
		return nil, nil, err
	}

	for _, p := range gas.imports {
		astutil.AddImport(fset, f, p)
	}

	gas.fset = fset
	actors := gas.actors
	halt := false
	n = astutil.Apply(f, func(c *astutil.Cursor) bool {
		for _, a := range actors {
			if a.Selector(c) {
				if e := a.Action(c); e != nil {
					err = e
					halt = true
				}
			}
		}

		return !halt
	}, nil)

	return
}

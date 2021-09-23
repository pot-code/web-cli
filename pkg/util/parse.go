package util

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/pot-code/web-cli/pkg/constants"
	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/ast/astutil"
)

type GoModMeta struct {
	Author      string
	ProjectName string
	Requires    []*modfile.Require
}

func ParseGoMod(path string) (*GoModMeta, error) {
	fd, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.WithMessage(err, "can't locate go.mod")
		}
		return nil, errors.Wrap(err, "failed to open go.mod")
	}

	content, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read")
	}

	mp, err := modfile.Parse(constants.GoModFile, content, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse mod file")
	}

	meta := new(GoModMeta)
	meta.Requires = mp.Require

	url := mp.Module.Mod.Path
	parts := strings.Split(url, "/")
	meta.Author = parts[len(parts)-2]
	meta.ProjectName = parts[len(parts)-1]

	if meta.Author == "" || meta.ProjectName == "" {
		return nil, errors.New("failed to extrace meta data from go.mod")
	}
	return meta, nil
}

type GoAstActor interface {
	//Selector select target node
	Selector(c *astutil.Cursor) bool
	// Action do action on selected node
	Action(c *astutil.Cursor) error
}

type GoFileParser struct {
	file    string
	fset    *token.FileSet
	actors  []GoAstActor
	imports []string
}

func NewGoFileParser(file string) *GoFileParser {
	return &GoFileParser{file: file, actors: []GoAstActor{}, imports: make([]string, 0)}
}

func (gas *GoFileParser) AddActor(actor ...GoAstActor) {
	gas.actors = append(gas.actors, actor...)
}

func (gas *GoFileParser) AddImport(p string) {
	gas.imports = append(gas.imports, p)
}

func (gas *GoFileParser) Parse() (fset *token.FileSet, n ast.Node, err error) {
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
					err = errors.WithStack(e)
					halt = true
				}
			}
		}

		return !halt
	}, nil)

	return
}

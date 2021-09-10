package util

import (
	"bufio"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/ast/astutil"
)

type GoModMeta struct {
	Author, ProjectName string
}

func ParseGoMod(path string) (*GoModMeta, error) {
	fd, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.WithMessage(err, "can't locate go.mod")
		}
		return nil, errors.Wrap(err, "failed to open go.mod")
	}

	var author, projectName string
	reader := bufio.NewScanner(fd)
	for reader.Scan() {
		line := reader.Text()
		if strings.HasPrefix(line, "module") {
			url := strings.TrimPrefix(line, "module ")
			parts := strings.Split(url, "/")
			author = parts[len(parts)-2]
			projectName = parts[len(parts)-1]
			break
		}
	}
	if author == "" || projectName == "" {
		return nil, errors.New("failed to extrace meta data from go.mod")
	}
	return &GoModMeta{author, projectName}, nil
}

type GoAstActor interface {
	//Selector select target node
	Selector(c *astutil.Cursor) bool
	// Action do action on selected node
	Action(c *astutil.Cursor) error
}

type GoFileParser struct {
	file   string
	fset   *token.FileSet
	actors []GoAstActor
}

func NewGoAstStore(file string) *GoFileParser {
	return &GoFileParser{file: file, actors: []GoAstActor{}}
}

func (gas *GoFileParser) AddActor(actor ...GoAstActor) {
	gas.actors = append(gas.actors, actor...)
}

func (gas *GoFileParser) Parse() (n ast.Node, err error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, gas.file, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	gas.fset = fset
	actors := gas.actors
	halt := false
	modifiedAst := astutil.Apply(f, func(c *astutil.Cursor) bool {
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

	return modifiedAst, nil
}

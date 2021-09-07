package cmd

import (
	"errors"
	"go/ast"
	"go/token"
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	ErrUnhandledEt             = errors.New("no handler for array element type")
	ErrUnhandledField          = errors.New("no handler for field type")
	ErrUnhandledSelectorPrefix = errors.New("no handler for selector prefix")
	ErrUnhandledPflagMapping   = errors.New("no handler for pflag mapping")
)

type pflagType struct {
	Type         string
	DefaultValue string
}

var pflagMap = map[string]pflagType{
	"string":        {"pflag.String", `""`},
	"int32":         {"pflag.Int32", "0"},
	"int":           {"pflag.Int", "0"},
	"[]string":      {"pflag.StringSlice", "[]string{}"},
	"[]int":         {"pflag.IntSlice", "[]int{}"},
	"time.Duration": {"pflag.Duration", "0"},
}

type pflagFieldDescriptor struct {
	FlagType string
	Help     string
	TagMeta  *tagMeta
}

type configStructVisitor struct {
	fset *token.FileSet
	fds  []*pflagFieldDescriptor
	pkg  string
}

func newConfigStructVisitor(fset *token.FileSet, pkg string) *configStructVisitor {
	return &configStructVisitor{
		fset: fset,
		fds:  []*pflagFieldDescriptor{},
		pkg:  pkg,
	}
}

func (cv *configStructVisitor) visitField(node ast.Node, prefix []string) {
	n := node.(*ast.Field)
	comment := n.Comment.Text()
	fieldName := n.Names[0].String()

	tag := n.Tag.Value
	ptm, err := parseTag(tag)
	if err != nil {
		log.WithFields(log.Fields{
			"fieldname": fieldName,
			"position":  cv.fset.Position(node.Pos()),
		}).Error(err)
		panic(err)
	}

	comment = strings.TrimSpace(comment)

	log.WithFields(log.Fields{
		"tag_raw":              tag,
		"tag_meta.Default":     ptm.Default,
		"tag_meta.MarshalName": ptm.MarshalName,
		"tag_meta.Required":    ptm.Required,
		"comment":              comment,
		"field_name":           fieldName,
	}).Debug("field info")

	switch nt := n.Type.(type) {
	case *ast.Ident:
		cv.visitIdent(nt, append(prefix, ptm.MarshalName), nt.String(), comment, ptm)
	case *ast.StructType:
		cv.visitStruct(nt, append(prefix, ptm.MarshalName))
	case *ast.SelectorExpr:
		cv.visitSelectorExpr(nt, append(prefix, ptm.MarshalName), comment, ptm)
	case *ast.ArrayType:
		cv.visitArray(nt, append(prefix, ptm.MarshalName), comment, ptm)
	default:
		log.WithFields(log.Fields{
			"fieldname": fieldName,
			"position":  cv.fset.Position(node.Pos()),
		}).Error(ErrUnhandledField)
		panic(ErrUnhandledField)
	}
}

func (cv *configStructVisitor) visitIdent(node ast.Node, prefix []string, ft, comment string, tm *tagMeta) {
	n := node.(*ast.Ident)

	if pt, ok := pflagMap[ft]; ok {
		if tm.Default == "" {
			log.WithFields(log.Fields{
				"position": cv.fset.Position(node.Pos()),
			}).Warn("no default value specified")
			tm.Default = pt.DefaultValue
		}

		tm.MarshalName = strings.Join(prefix, ".")
		cv.fds = append(cv.fds, &pflagFieldDescriptor{
			FlagType: pt.Type,
			Help:     comment,
			TagMeta:  tm,
		})
	} else {
		log.WithFields(log.Fields{
			"go_type":  ft,
			"position": cv.fset.Position(n.Pos()),
		}).Error(ErrUnhandledPflagMapping)
		panic(ErrUnhandledPflagMapping)
	}
}

func (cv *configStructVisitor) visitSelectorExpr(node ast.Node, prefix []string, comment string, tm *tagMeta) {
	n := node.(*ast.SelectorExpr)

	switch xt := n.X.(type) {
	case *ast.Ident:
		ft := xt.String() + "." + n.Sel.String()
		cv.visitIdent(xt, prefix, ft, comment, tm)
	default:
		log.WithFields(log.Fields{
			"position": cv.fset.Position(node.Pos()),
		}).Error(ErrUnhandledSelectorPrefix)
		panic(ErrUnhandledSelectorPrefix)
	}
}

func (cv *configStructVisitor) visitArray(node ast.Node, prefix []string, comment string, tm *tagMeta) {
	n := node.(*ast.ArrayType)

	switch nt := n.Elt.(type) {
	case *ast.Ident:
		elt := nt.String()
		ft := "[]" + elt
		cv.visitIdent(nt, prefix, ft, comment, tm)
	default:
		log.WithFields(log.Fields{
			"position": cv.fset.Position(node.Pos()),
		}).Error(ErrUnhandledEt)
		panic(ErrUnhandledEt)
	}
}

func (cv *configStructVisitor) visitStruct(node ast.Node, prefix []string) {
	n := node.(*ast.StructType)

	for _, f := range n.Fields.List {
		cv.visitField(f, prefix)
	}
}

func (cv *configStructVisitor) visit(node ast.Node) {
	n := node.(*ast.StructType)

	for _, f := range n.Fields.List {
		cv.visitField(f, []string{})
	}
}

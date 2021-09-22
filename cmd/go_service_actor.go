package cmd

import (
	"go/ast"

	"golang.org/x/tools/go/ast/astutil"
)

type AddWireSetActor struct {
	pkg string
}

func NewAddWireSetActor(pkg string) *AddWireSetActor {
	return &AddWireSetActor{pkg: pkg}
}

//Selector select target node
func (ah *AddWireSetActor) Selector(c *astutil.Cursor) bool {
	n := c.Node()

	vs, ok := n.(*ast.ValueSpec)
	if !ok {
		return false
	}

	return vs.Names[0].Name == "HttpSet"
}

// Action do action on selected node
func (ah *AddWireSetActor) Action(c *astutil.Cursor) error {
	n := c.Node()
	vs, _ := n.(*ast.ValueSpec)

	ce := vs.Values[0].(*ast.CallExpr)
	ce.Args = append(ce.Args, &ast.SelectorExpr{
		X:   &ast.Ident{Name: ah.pkg},
		Sel: &ast.Ident{Name: "HttpSet"},
	})
	return nil
}

// type AddHandlerActor struct {
// 	pkg     string
// 	handler string
// }

// //Selector select target node
// func (ah *AddHandlerActor) Selector(c *astutil.Cursor) bool {
// 	n := c.Node()

// 	ts, ok := n.(*ast.TypeSpec)
// 	if !ok {
// 		return false
// 	}

// 	return ts.Name.Name == "HandlerCollection"
// }

// // Action do action on selected node
// func (ah *AddHandlerActor) Action(c *astutil.Cursor) error {
// 	n := c.Node()
// 	ts, _ := n.(*ast.TypeSpec)

// 	st := ts.Type.(*ast.StructType)
// 	fields := st.Fields.List
// 	fields = append(fields, &ast.Field{
// 		Type: &ast.StarExpr{
// 			X: &ast.SelectorExpr{
// 				X:   &ast.Ident{Name: "port"},
// 				Sel: &ast.Ident{Name: ah.handler},
// 			},
// 		},
// 	})
// 	return nil
// }

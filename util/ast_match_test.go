package util

import "testing"
import "go/ast"
import "go/parser"
import "go/token"


const pattern = `
package foo

type myInterface interface {
	Read() string
}

// _METHOD_NAME is part of the _INTERFACE_NAME interface.
func (x *_STRUCT_NAME) _METHOD_NAME() _SLOT_TYPE {
	return x._SLOT_NAME
}
`  // End

func TestASTMatch(t *testing.T) {
	fset := token.NewFileSet()
	parsed, err := parser.ParseFile(fset, "pattern", pattern, parser.ParseComments)
	if err != nil {
		t.Fatalf("Error while parsing pattern: %s", err)
		return
	}
	show := func(what string, x interface{}) {
		t.Logf("%s: %T %v", what, x, x)
	}
	my_interface := parsed.Decls[0].(*ast.GenDecl).Specs[0].(*ast.TypeSpec)
	method := parsed.Decls[1].(*ast.FuncDecl)
	show("my_interface", my_interface)
	show("method", method)
	expectType := method.Type
	gotType := my_interface.Type.(*ast.InterfaceType).Methods.List[0].Type
	show("method FuncType", expectType)
 	show("my_interface 1st method Type", gotType)
	scratchpad :=  map[string]interface{}{}
	matched, err := AstMatch(expectType, gotType, scratchpad)
	t.Logf("AstMatch result %v %v %#v", matched, err, scratchpad)
	if !matched {
		t.Errorf("Didn't match: %v, %v", expectType, gotType)
	}
}


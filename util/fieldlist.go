package util

import "go/ast"

// FieldListSlice returns a []*ast.Field of the Fields in l.
// This is a convenience to deal with the cases where a *ast.FieldList
// valued slot contains nil rather than a FieldList with an empty
//  List slice.
func FieldListSlice(l *ast.FieldList) []*ast.Field {
	if l == nil {
		return []*ast.Field{}
	}
	return l.List
}

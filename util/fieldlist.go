package util

import "fmt"
import "path"
import "strconv"
import "strings"
import "go/ast"
import "go/types"


// FieldListSlice returns a []*ast.Field of the Fields in l.
// This is a convenience to deal with the cases where a *ast.FieldList
// valued slot contains nil rather than a FieldList with an empty
//  List slice.
func FieldListSlice(l *ast.FieldList) []*ast.Field {
	if l == nil || l.List == nil {
		return []*ast.Field{}
	}
	return l.List
}


// FieldListString returns a string representation of a FieldList
// suitable for including in code.
//
// If need_names is true then generated parameter names will be
// included in the formal list and a second string that is just a
// comma separated lst of those names is returned.
//
// If for_result is true then the resulting string will be wrapped in
// parentheses.
func FieldListString(l *ast.FieldList, info *types.Info, qualifier types.Qualifier, need_names bool, for_result bool) (string, string) {
	formal := []string{}
	names := []string{}
	for i, field := range FieldListSlice(l) {
		typestring := types.TypeString(info.Types[field.Type].Type, qualifier)
		if need_names {
			name := fmt.Sprintf("a%d", i)
			formal = append(formal, fmt.Sprintf("%s %s", name, typestring))
			names = append(names, name)
		} else {
			formal = append(formal, typestring)
		} 
	}
	s := strings.Join(formal, ", ")
	if for_result && s != "" {
		s = "(" + s + ")"
	}
	return s, strings.Join(names, ", ")
}


// TypeStringQualifier returns a types.Qualifier suitable for
// including the result of types.TypeString in code.
func TypeStringQualifier(f *ast.File) types.Qualifier {
	return func(pkg *types.Package) string {
		ppath := pkg.Path()
		for _, ispec := range f.Imports {
			unq, err := strconv.Unquote(ispec.Path.Value)
			if err != nil {
				panic(err)
			}
			if unq == ppath {
				if ispec.Name == nil {
					_, base := path.Split(ppath)
					/*
					f.PendingImports = append(f.PendingImports, &ToImport{
						Name: "",
						Path: ppath,
					})
					*/
					return base
				} else {
					/*
					f.PendingImports = append(f.PendingImports, &ToImport{
						Name: ispec.Name.Name,
						Path: ppath,
					})
					*/
					return ispec.Name.Name
				}
			}
		}
		return ""

	}
}

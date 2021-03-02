package util

import "strconv"
import "path/filepath"
import "go/ast"
import "go/token"
import "golang.org/x/tools/go/ast/astutil"


// ImportAdder is a function returned by ImportSpecMatch if the
// ImportSpec matches the specified name, and nil otherwise.  The
// ImportAdder can be called on a token.FileSet and an ast.File to add
// the ImportSpec to that file.  It will use astutil.AddImport or
// astutil.AddNamedImport to add the import.
type ImportAdder func(fset *token.FileSet, file *ast.File)


// ImportSpecMatch returns an ImportAdder if the ImportSpec matches the
// specified name.  It can return an error if the Path of the
// ImportSpec isn't a quoted expression
func ImportSpecMatch(ispec *ast.ImportSpec, name string) (ImportAdder, error) {
	unq, err := strconv.Unquote(ispec.Path.Value)
	if err != nil {
		return ImportAdder(nil), err
	}
	if ispec.Name == nil {
		if filepath.Base(unq) == name {
			return ImportAdder(func(fset *token.FileSet, file *ast.File) {
				astutil.AddImport(fset, file, unq)
			}), nil
		}
	} else {
		if  ispec.Name.Name == name {
			return ImportAdder(func(fset *token.FileSet, file *ast.File) {
				astutil.AddNamedImport(fset, file, name, unq)
			}), nil
		}
	}
	return ImportAdder(nil), err
}


type visitor struct {
	fset *token.FileSet
	in *ast.File
	out *ast.File
	errors []error
}


func LeftmostSelector(sel *ast.SelectorExpr) *ast.Ident {
	switch e := sel.X.(type) {
	case *ast.Ident:
		return e
	case *ast.SelectorExpr:
		return LeftmostSelector(e)
	default:
		return nil
	}
}


func (v *visitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.SelectorExpr:
		left := LeftmostSelector(n)
		if left != nil {
			// left might be a package reference
			for _, ispec := range v.in.Imports {
				f, err := ImportSpecMatch(ispec, left.Name)
				if err != nil {
					v.errors = append(v.errors, err)
					continue
				}
				if f != nil {
					f(v.fset, v.out)
					break
				}
			}
		}
	}
	return v
}


// EnsureImports makes sure that out has all imports needed to satisfy
// any references.
func EnsureImports(fset *token.FileSet, in, out *ast.File) []error {
	v := &visitor{
		fset: fset,
		in: in,
		out: out,
	}
	ast.Walk(v, out)
	return v.errors
}


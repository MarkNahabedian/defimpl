// Looking up interface definitions.
package main

import "fmt"
import "go/ast"


// IDKey (interface definition key) is used to identify the
// InterfaceDefinition to be found.
type IDKey struct {
	Package string
	Name string
}

func (key *IDKey) String() string {
	return fmt.Sprintf("%s.%s", key.Package, key.Name)
}

func ExprToIDKey(e ast.Expr, defaultPkg string) *IDKey {
	var etk func (recursive bool, e ast.Expr, defaultPkg string) *IDKey
	etk = func(recursive bool, e ast.Expr, defaultPkg string) *IDKey {
		switch e1 := e.(type) {
		case *ast.Ident:
			return &IDKey{
				Package: defaultPkg,
				Name: e1.Name,
			}
		case *ast.SelectorExpr:
			if recursive {
				panic("Nexted selector expressions not supported")
			}
			p, ok := e1.X.(*ast.Ident)
			if !ok {
				panic(fmt.Sprintf("Unsupported selector %#v", e1))
			}
			return etk(true, e1.Sel, p.Name)
		default:
			panic(fmt.Sprintf("Unsupported expression type %T", e1))
		}
	}
	return etk(false, e, defaultPkg)
}

func (ctx *context) IDLookup(key *IDKey) *InterfaceDefinition {
	for _, f := range ctx.files {
		found := f.IDLookup(key)
		if found != nil {
			return found
		}		
	}
	return nil
}

func (f *File) IDLookup(key *IDKey) *InterfaceDefinition {
	for _, i := range f.Interfaces {
		found := i.IDLookup(key)
		if found != nil {
			return found
		}
	}
	return nil
}

func (i *InterfaceDefinition) IDLookup(key *IDKey) *InterfaceDefinition {
	if i.InterfaceName != key.Name {
		return nil
	}
	if i.Package != key.Package {
		return nil
	}
	return i
}
	

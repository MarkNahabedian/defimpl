package main

import "defimpl/util"
import "fmt"
import "go/ast"
import "strings"


type InterfaceDefinition struct {
	File          *File
	IsAbstract    bool
	InterfaceType *ast.InterfaceType
	InterfaceName string
	VerbPhrases   []VerbPhrase
	Inherited     []*IDKey                // Interfaces that are included by this one
	AllInherited  []*InterfaceDefinition  // Transitive closure of all inherited interfaces.
}

func (idef *InterfaceDefinition) QualifiedName() string {
	return util.ImplName(idef.Package(), idef.InterfaceName)
}

// StructName returns the default name for the implementing struct for
// the interface represented by this InterfaceDefinition.
func (idef *InterfaceDefinition) StructName() string {
	return idef.InterfaceName + "Impl"
}

// DefinesStruct returns true if an implementing sruct should be
// defined for the interface represented by this InterfaceDefinition.
func (idef *InterfaceDefinition) DefinesStruct() bool {
	return !idef.IsAbstract
}

func (idef *InterfaceDefinition) Fields() []*ast.Field {
	return util.FieldListSlice(idef.InterfaceType.Methods)
}

func (idef *InterfaceDefinition) Package() string {
	return idef.File.Package
}


const InterfaceIsAbstractMarker string = "(ABSTRACT)"

// isAbstractInterface returns true if the declaration -- which should
// define an interface -- has a comment with the "abstract" token.
func isAbstractInterface(x *ast.GenDecl) bool {
	var hasAbstract = func(cmnt *ast.CommentGroup) bool {
		if cmnt == nil || cmnt.List == nil {
			return false
		}
		for _, c := range cmnt.List {
			if strings.Contains(c.Text, InterfaceIsAbstractMarker) {
				return true
			}
		}
		return false
	}
	return hasAbstract(x.Doc)
}


// NewInterface returns a new InterfaceDefinition if decl represents
// an interface definition, otherwise it returns nil.
func NewInterface(ctx *context, file *File, decl ast.Decl) *InterfaceDefinition {
	gd, ok := decl.(*ast.GenDecl)
	if !ok {
		return nil
	}
	spec, ok := gd.Specs[0].(*ast.TypeSpec)
	if !ok {
		return nil
	}
	it, ok := spec.Type.(*ast.InterfaceType)
	if !ok {
		return nil
	}
	if len(gd.Specs) > 1 {
		// Apparently I have insufficient understanding of what interface type specs look like.
		panic(fmt.Sprintf("type definition of an interface type has more than one Spec: %s\n",
			ctx.fset.Position(gd.TokPos).String()))
	}
	id := &InterfaceDefinition{
		File: file,
		// It appears that the parser associates the comment group with
		// the outer GenDecl rather than with the TypeSpec.
		IsAbstract:    isAbstractInterface(gd),
		InterfaceType: it,
		InterfaceName: spec.Name.Name,
		Inherited:     []*IDKey{},
	}
	for _, m := range id.Fields() {
		GetVerbPhrase(ctx, id, m)
	}
	return id
}


// TypePackage returns the package name if the Expr (which should
// identify a type) specifies one.
func TypePackage(t ast.Expr) string {
	if t == nil {
		return ""
	}
	var tp func(ast.Expr, bool) string
	tp = func(t ast.Expr, top bool) string {
		switch e := t.(type) {
		case *ast.Ident:
			if top {
				return ""
			}
			return e.Name
		case *ast.SelectorExpr:
			return tp(e.X, false)
		case *ast.ArrayType:
			return tp(e.Elt, true)
		case *ast.StarExpr:
			return tp(e.X, true)
		case *ast.FuncType:
			// Unnamed function, so no package.
			return ""
		default:
			panic(fmt.Sprintf("TypePackage: unsupported expression type %T", t))
		}
	}
	return tp(t, true)
}


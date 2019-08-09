package main

import "defimpl/util"
import "fmt"
import "go/ast"
import "go/token"
import "os"
import "path/filepath"
import "reflect"
import "strconv"
import "strings"
import "golang.org/x/tools/go/ast/astutil"

type InterfaceDefinition struct {
	IsAbstract    bool
	InterfaceType *ast.InterfaceType
	InterfaceName string
	SlotSpecs     []*slotSpec
	Package       string
	Inherited     []*IDKey                // Interfaces that are included by this one
	AllInherited  []*InterfaceDefinition  // Transitive closure of all inherited interfaces.
}

func (idef *InterfaceDefinition) QualifiedName() string {
	return util.ImplName(idef.Package, idef.InterfaceName)
}

func (idef *InterfaceDefinition) StructName() string {
	return idef.InterfaceName + "Impl"
}

func (idef *InterfaceDefinition) AddImports(ctx *context, in, out *ast.File) {
	for _, ss := range idef.SlotSpecs {
		ss.AddImports(ctx, in, out)
	}
}

func (idef *InterfaceDefinition) DefinesStruct() bool {
	return !idef.IsAbstract && (len(idef.SlotSpecs) > 0)
}

const InterfaceIsAbstractMarker string = "(ABSTRACT)"

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
func NewInterface(ctx *context, pkg string, decl ast.Decl) *InterfaceDefinition {
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
		// It appears that the parser associates to comment group with
		// the outer GenDecl rather than with the TypeSpec.
		IsAbstract:    isAbstractInterface(gd),
		InterfaceType: it,
		InterfaceName: spec.Name.Name,
		SlotSpecs:     []*slotSpec{},
		Package:       pkg,
		Inherited:     []*IDKey{},
	}
	getSpec := func(name string) *slotSpec {
		for _, sspec := range id.SlotSpecs {
			if sspec.Name == name {
				return sspec
			}
		}
		spec := &slotSpec{Name: name}
		id.SlotSpecs = append(id.SlotSpecs, spec)
		return spec
	}
	for _, m := range util.FieldListSlice(it.Methods) {
		if len(m.Names) == 0 {
			id.Inherited = append(id.Inherited, ExprToIDKey(m.Type, id.Package))
		}
		if len(m.Names) != 1 {
			continue
		}
		verb, slot := methodDefimpl(m)
		if verb == nil {
			continue
		}
		spec := getSpec(slot)
		spec.assimilate(ctx, id, m, verb)
	}
	return id
}

// methodDefimpl returns the verb and slot name from a method's
// defimpl comment if any.
func methodDefimpl(method *ast.Field) (verb *VerbDefinition, slot_name string) {
	if method.Comment == nil {
		return nil, ""
	}
	for _, c := range method.Comment.List {
		val, ok := reflect.StructTag(c.Text[2:]).Lookup("defimpl")
		if !ok {
			continue
		}
		split := strings.Split(val, " ")
		vd := LookupVerb(split[0])
		if vd != nil {
			return vd, split[1]
		}
	}
	return nil, ""
}

type slotSpec struct {
	Name  string
	Type  ast.Expr // types.Type
	Verbs []*VerbTemplateParameter
}

// CheckType fills in the Type of spec and makes sure that the Type is
// consistent for all method verbs associated with that slot.
func (spec *slotSpec) CheckType(typ ast.Expr) error {
	teq := func(t1, t2 ast.Expr) bool {
		return reflect.TypeOf(t1) == reflect.TypeOf(t2)
	}
	if spec.Type == nil {
		spec.Type = typ
	} else {
		if !teq(spec.Type, typ) {
			return fmt.Errorf("Incompatible types: %#v, %#v",
				typ, spec.Type)
		}
	}
	return nil
}

func (spec *slotSpec) assimilate(ctx *context, id *InterfaceDefinition, m *ast.Field, verb *VerbDefinition) {
	errorPosition := func() token.Position {
		return ctx.fset.Position(m.Comment.Pos())
	}
	if err := verb.Assimilate(ctx, verb, spec, id, m); err != nil {
		fmt.Fprintf(os.Stderr, "%s at %s\n", err, errorPosition().String())
		return
	}
}

func (spec *slotSpec) AddImports(ctx *context, in, out *ast.File) {
	p := TypePackage(spec.Type)
	if p == "" {
		return
	}
	AddImport(ctx.fset, in, out, p)
}

// TypePackage returns the package name if the Expr (which should
// identify a type) specifies one.
func TypePackage(t ast.Expr) string {
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
		default:
			panic(fmt.Sprintf("TypePackage: unsupported expression type %T", t))
		}
	}
	return tp(t, true)
}

func AddImport(fset *token.FileSet, in, out *ast.File, pkg_identifier string) {
	for _, para := range astutil.Imports(fset, in) {
		for _, ispec := range para {
			unq, err := strconv.Unquote(ispec.Path.Value)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				continue
			}
			if ispec.Name == nil {
				if filepath.Base(unq) == pkg_identifier {
					astutil.AddImport(fset, out, unq)
					return
				}
			} else {
				if  ispec.Name.Name == pkg_identifier {
					astutil.AddNamedImport(fset, out, pkg_identifier, unq)
					return
				}
			}
		}
	}
	panic(fmt.Sprintf("Can't find package %s in %s",
		pkg_identifier, fset.Position(in.Package).Filename))
}

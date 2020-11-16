package main

import "defimpl/util"
import "fmt"
import "go/ast"
import "go/token"
import "os"
import "reflect"
import "strings"


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

func (idef *InterfaceDefinition) DefinesStruct() bool {
	return !idef.IsAbstract && (len(idef.SlotSpecs) > 0)
}

func (idef *InterfaceDefinition) Fields() []*ast.Field {
	return util.FieldListSlice(idef.InterfaceType.Methods)
}

func (idef *InterfaceDefinition) GetSpec(name string) *slotSpec {
	for _, sspec := range id.SlotSpecs {
		if sspec.Name == name {
			return sspec
		}
	}
	spec := &slotSpec{Name: name}
	id.SlotSpecs = append(id.SlotSpecs, spec)
	return spec
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
		// It appears that the parser associates the comment group with
		// the outer GenDecl rather than with the TypeSpec.
		IsAbstract:    isAbstractInterface(gd),
		InterfaceType: it,
		InterfaceName: spec.Name.Name,
		SlotSpecs:     []*slotSpec{},
		Package:       pkg,
		Inherited:     []*IDKey{},
	}
	for _, m := range id.Fields() {
		if len(m.Names) == 0 {
			id.Inherited = append(id.Inherited, ExprToIDKey(m.Type, id.Package))
		}
		if len(m.Names) != 1 {
			continue
		}
		verb, slot := methodDefimpl(ctx, m)
		if verb == nil {
			continue
		}
		spec := id.GetSpec(slot)
		spec.assimilate(ctx, id, m, verb)
	}
	return id
}

// methodDefimpl returns the verb and slot name from a method's
// defimpl comment if any.
func methodDefimpl(ctx *context, method *ast.Field) (verb *VerbDefinition, slot_name string) {
	if method.Comment == nil {
		return nil, ""
	}
	for _, c := range method.Comment.List {
		val, ok := reflect.StructTag(c.Text[2:]).Lookup("defimpl")
		if !ok {
			continue
		}
		pos := ctx.fset.Position(c.Slash)
		split := strings.Split(val, " ")
		if len(split) < 1 {
			fmt.Fprintf(os.Stderr, "defimpl: No verb in defimpl comment %s: %q\n", pos, c.Text)
			continue
		}
		verb := split[0]
		vd := LookupVerb(verb)
		if vd == nil {
			fmt.Fprintf(os.Stderr, "defimpl: Unknown verb %q in defimpl comment %s: %q\n", verb, pos, c.Text)
			continue			
		}
		if len(split) != vd.ParamCount + 1 {
			fmt.Fprintf(os.Stderr, "defimpl verb %q expects %d parameters: %s %q",
				verb, vd.ParamCount, pos, c.Text)
			continue
		}
		slot := ""
		if len(split) > 1 {
			slot = split[1]
		}
		if verbose {
			fmt.Printf("%s: %q %q %v\n", pos, verb, slot, vd != nil)
		}
		if vd != nil {
			return vd, slot
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
		fmt.Fprintf(os.Stderr, "defimpl: %s at %s\n", err, errorPosition().String())
		return
	}
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


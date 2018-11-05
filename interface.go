package main

import "defimpl/util"
import "fmt"
import "go/ast"
import "go/token"
// import "go/types"
import "os"
import "reflect"
import "strings"
import _ "golang.org/x/tools/go/ast/astutil"

type InterfaceDefinition struct {
	InterfaceType *ast.InterfaceType
	InterfaceName string
	SlotSpecs []*slotSpec
}

func (idef *InterfaceDefinition) StructName() string {
	return idef.InterfaceName + "Impl"
}

func (idef *InterfaceDefinition) AddImports(ctx *context, in, out *ast.File) {
	for _, ss := range idef.SlotSpecs {
		ss.AddImports(ctx, in, out)
	}
}

// NewInterface returns a new InterfaceDefinition if decl represents
// an interface definition, otherwise it returns nil.
func NewInterface(ctx *context, decl ast.Decl) *InterfaceDefinition {
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
		InterfaceType: it,
		InterfaceName: spec.Name.Name,
		SlotSpecs: []*slotSpec{},
	}
	specs := map[string] *slotSpec{}
	getSpec := func(name string) *slotSpec {
		if spec, ok := specs[name]; ok {
			return spec
		}
		spec := &slotSpec{ Name: name }
		specs[name] = spec
		return spec
	}
	for _, m := range util.FieldListSlice(it.Methods) {
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
	for _, sspec := range specs {
		id.SlotSpecs = append(id.SlotSpecs, sspec)
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
	Name string
	Type ast.Expr    // types.Type
	Verbs []*VerbTemplateParameter
}

// CheckType fills in the Type of spec and makes sure that the Type is
// consistent for all method verbs associated with that slot.
func (spec *slotSpec) CheckType(typ ast.Expr) error {
	teq := func (t1, t2 ast.Expr) bool {
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
/*
	typestring := spec.Type.Value
	split := strings.Split(typestring, ".")
	if len(split) <= 1 {
		return
	}
	// *** Need to split by / and take last element
	if strings.HasSuffix(split[0], out.Name.Name) {
		return
	}
	// *** need to lookup package name in in to get the path and make
	// sure the same name is defined as that path in out.
	astutil.AddImport(ctx.fset, out, split[0])
*/
}


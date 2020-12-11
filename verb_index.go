package main

import "go/ast"
import "go/types"
import "text/template"


type IndexVerbPhrase struct {
	slotVerbPhrase
}

var _ VerbPhrase = (*IndexVerbPhrase)(nil)
var _ SlotVerbPhrase = (*IndexVerbPhrase)(nil)
var _ MethodTemplateParameter = (*IndexVerbPhrase)(nil)


type Verb_Index struct {
	slotVerbDefinition
}

var _ VerbDefinition = (*Verb_Index)(nil)

func init() {
	vd := &Verb_Index{}
	VerbDefinitions[vd.Tag()] = vd
}

// Verb is part of the VerbDefinition interface.
func (vd *Verb_Index) Tag() string { return "index" }

// Description is part of the VerbDefinition interface.
func (vd *Verb_Index) Description() string {
	return "returns the element of the specified slice valued field at the specified (zero based) index."
}

// NewVerbPhrase is part of the VerbDefinition interface.
func (vd *Verb_Index) NewVerbPhrase(ctx *context, idef *InterfaceDefinition, field *ast.Field, comment *ast.Comment) (VerbPhrase, error) {
	slot, err := parse_slot_verb_phrase(ctx, field, comment)
		if err != nil {
		return nil, err
	}
	slot_type, err := CheckSignatures(ctx, vd, idef.Package(), field, vd.MethodTemplate())
	if err != nil {
		return nil, err
	}
	vp := &IndexVerbPhrase{
		slotVerbPhrase {
			baseVerbPhrase: baseVerbPhrase {
				verb: vd,
				idef: idef,
				field: field,
			},
			slot_name: slot,
			slot_type: types.NewSlice(slot_type),
		},
	}
	if err := addSlotSpec(idef, vp); err != nil {
		return nil, err
	}
	return vp, nil
}

var index_method_template = template.Must(
	template.New("index_method_template").Parse(`
// {{.MethodName}} is part of the {{.InterfaceName}} interface.
func (x *{{.StructName}}) {{.MethodName}} (index int) {{.TypeString .SlotType.Elem}} {
	return x.{{.SlotName}}[index]
}
`))

// MethodTemplate is part of the VerbDefinition interface.
func (vd *Verb_Index) MethodTemplate() *template.Template {
	return index_method_template
}


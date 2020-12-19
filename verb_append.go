package main

import "go/ast"
import "go/types"
import "text/template"


type AppendVerbPhrase struct {
	slotVerbPhrase
}

var _ VerbPhrase = (*AppendVerbPhrase)(nil)
var _ SlotVerbPhrase = (*AppendVerbPhrase)(nil)
var _ GlobalsTemplateParameter = (*AppendVerbPhrase)(nil)


type Verb_Append struct {
	slotVerbDefinition
}

var _ VerbDefinition = (*Verb_Append)(nil)

func init() {
	vd := &Verb_Append{}
	VerbDefinitions[vd.Tag()] = vd
}

// Verb is part of the VerbDefinition interface.
func (vd *Verb_Append) Tag() string { return "append" }

// Description is part of the VerbDefinition interface.
func (vd *Verb_Append) Description() string {
	return "appends the specified values to the field."
}

// NewVerbPhrase is part of the VerbDefinition interface.
func (vd *Verb_Append) NewVerbPhrase(ctx *context, idef *InterfaceDefinition, field *ast.Field, comment *ast.Comment) (VerbPhrase, error) {
	slot, err := parse_slot_verb_phrase(ctx, field, comment)
	if err != nil {
		return nil, err
	}
	slot_type, err := CheckSignatures(ctx, vd, idef.Package(), field, vd.GlobalsTemplate())
	if err != nil {
		return nil, err
	}
	vp := &AppendVerbPhrase{
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

var append_method_template =  template.Must(
	template.New("append_method_template").Parse(`
// {{.MethodName}} is part of the {{.InterfaceName}} interface.
func (x *{{.StructName}}) {{.MethodName}} (v ...{{.TypeString .SlotType.Elem}}) {
	x.{{.SlotName}} = append(x.{{.SlotName}}, v...)
}
`))

// GlobalsTemplate is part of the VerbDefinition interface.
func (vd *Verb_Append) GlobalsTemplate() *template.Template {
	return append_method_template
}


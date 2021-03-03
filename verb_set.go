package main

import "go/ast"
import "text/template"


type SetVerbPhrase struct {
	slotVerbPhrase
}

var _ VerbPhrase = (*SetVerbPhrase)(nil)
var _ SlotVerbPhrase = (*SetVerbPhrase)(nil)
var _ GlobalsTemplateParameter = (*SetVerbPhrase)(nil)


type Verb_Set struct {
	slotVerbDefinition
}

var _ VerbDefinition = (*Verb_Set)(nil)

func init() {
	vd := &Verb_Set{}
	VerbDefinitions[vd.Tag()] = vd
}

// Verb is part of the VerbDefinition interface.
func (vd *Verb_Set) Tag() string { return "set" }

// Description is part of the VerbDefinition interface.
func (vd *Verb_Set) Description() string {
	return "sets the value of the field to that provided."
}

// NewVerbPhrase is part of the VerbDefinition interface.
func (vd *Verb_Set) NewVerbPhrase(ctx *context, idef *InterfaceDefinition, field *ast.Field, comment *ast.Comment) (VerbPhrase, error) {
	slot, err := parse_slot_verb_phrase(ctx, field, comment)
	if err != nil {
		return nil, err
	}
	slot_type, err, _ := CheckSignatures(ctx, vd, idef.Package(), field, vd.GlobalsTemplate())
	if err != nil {
		return nil, err
	}
	vp := &SetVerbPhrase{
		slotVerbPhrase {
			baseVerbPhrase: baseVerbPhrase {
				verb: vd,
				idef: idef,
				field: field,
			},
			slot_name: slot,
			slot_type: slot_type,
		},
	}
	addSlotSpec(idef, vp)
	return vp, nil
}

var set_method_template = template.Must(
	template.New("set_method_template").Parse(`
// {{.MethodName}} is part of the {{.InterfaceName}} interface.  defimpl verb {{.Verb.Tag}}.
func (x *{{.StructName}}) {{.MethodName}}(v {{.TypeString .SlotType}}) {
	x.{{.SlotName}} = v
}
`))

// GlobalsTemplate is part of the VerbDefinition interface.
func (vd *Verb_Set) GlobalsTemplate() *template.Template {
	return set_method_template
}


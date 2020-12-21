package main

import "go/ast"
import "text/template"


type ReadVerbPhrase struct {
	slotVerbPhrase
}

var _ VerbPhrase = (*ReadVerbPhrase)(nil)
var _ SlotVerbPhrase = (*ReadVerbPhrase)(nil)
var _ GlobalsTemplateParameter = (*ReadVerbPhrase)(nil)


type Verb_Read struct {
	slotVerbDefinition
}

var _ VerbDefinition = (*Verb_Read)(nil)

func init() {
	vd := &Verb_Read{}
	VerbDefinitions[vd.Tag()] = vd
}

// Verb is part of the VerbDefinition interface.
func (vd *Verb_Read) Tag() string { return "read" }

// Description is part of the VerbDefinition interface.
func (vd *Verb_Read) Description() string {
	return "returns the value of the field."
}

// NewVerbPhrase is part of the VerbDefinition interface.
func (vd *Verb_Read) NewVerbPhrase(ctx *context, idef *InterfaceDefinition, field *ast.Field, comment *ast.Comment) (VerbPhrase, error) {
	slot, err := parse_slot_verb_phrase(ctx, field, comment)
	if err != nil {
		return nil, err
	}
	slot_type, err := CheckSignatures(ctx, vd, idef.Package(), field, vd.GlobalsTemplate())
	if err != nil {
		return nil, err
	}
	vp := &ReadVerbPhrase{
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

var read_method_template = template.Must(
	template.New("read_method_template").Parse(`
// {{.MethodName}} is part of the {{.InterfaceName}} interface.  defimpl verb {{.Verb.Tag}}.
func (x *{{.StructName}}) {{.MethodName}}() {{.TypeString .SlotType}} {
	return x.{{.SlotName}}
}
`))

// GlobalsTemplate is part of the VerbDefinition interface.
func (vd *Verb_Read) GlobalsTemplate() *template.Template {
	return read_method_template
}


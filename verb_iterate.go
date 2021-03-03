package main

import "go/ast"
import "go/types"
import "text/template"


type IterateVerbPhrase struct {
	slotVerbPhrase
}

var _ VerbPhrase = (*IterateVerbPhrase)(nil)
var _ SlotVerbPhrase = (*IterateVerbPhrase)(nil)
var _ GlobalsTemplateParameter = (*IterateVerbPhrase)(nil)


type Verb_Iterate struct {
	slotVerbDefinition
}

var _ VerbDefinition = (*Verb_Iterate)(nil)

func init() {
	vd := &Verb_Iterate{}
	VerbDefinitions[vd.Tag()] = vd
}

// Verb is part of the VerbDefinition interface.
func (vd *Verb_Iterate) Tag() string { return "iterate" }

// Description is part of the VerbDefinition interface.
func (vd *Verb_Iterate) Description() string {
	return "applies the specified function to each element of the slice-valued slot until the function returns false."
}

// NewVerbPhrase is part of the VerbDefinition interface.
func (vd *Verb_Iterate) NewVerbPhrase(ctx *context, idef *InterfaceDefinition, field *ast.Field, comment *ast.Comment) (VerbPhrase, error) {
	slot, err := parse_slot_verb_phrase(ctx, field, comment)
	if err != nil {
		return nil, err
	}
	slot_type, err, _ := CheckSignatures(ctx, vd, idef.Package(), field, vd.GlobalsTemplate())
	if err != nil {
		return nil, err
	}
	vp := &IterateVerbPhrase{
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

var iterate_method_template = template.Must(
	template.New("iterate_method_template").Parse(`
// {{.MethodName}} is part of the {{.InterfaceName}} interface.  defimpl verb {{.Verb.Tag}}.
func (x *{{.StructName}}) {{.MethodName}} (f func(item {{.TypeString .SlotType.Elem}}) bool) {
	for _, v := range x.{{.SlotName}} {
		if !f(v) {
			break
		}
	}
}
`))

// GlobalsTemplate is part of the VerbDefinition interface.
func (vd *Verb_Iterate) GlobalsTemplate() *template.Template {
	return iterate_method_template
}


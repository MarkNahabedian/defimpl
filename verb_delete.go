package main

import "go/ast"
import "go/types"
import "text/template"


type DeleteVerbPhrase struct {
	slotVerbPhrase
}

var _ VerbPhrase = (*DeleteVerbPhrase)(nil)
var _ SlotVerbPhrase = (*DeleteVerbPhrase)(nil)
var _ GlobalsTemplateParameter = (*DeleteVerbPhrase)(nil)


type Verb_Delete struct {
	slotVerbDefinition
}

var _ VerbDefinition = (*Verb_Delete)(nil)

func init() {
	vd := &Verb_Delete{}
	VerbDefinitions[vd.Tag()] = vd
}

// Verb is part of the VerbDefinition interface.
func (vd *Verb_Delete) Tag() string { return "delete" }

// Description is part of the VerbDefinition interface.
func (vd *Verb_Delete) Description() string {
	return "deletes the specified item from the filed."
}

// NewVerbPhrase is part of the VerbDefinition interface.
func (vd *Verb_Delete) NewVerbPhrase(ctx *context, idef *InterfaceDefinition, field *ast.Field, comment *ast.Comment) (VerbPhrase, error) {
	slot, err := parse_slot_verb_phrase(ctx, field, comment)
	if err != nil {
		return nil, err
	}
	slot_type, err := CheckSignatures(ctx, vd, idef.Package(), field, vd.GlobalsTemplate())
	if err != nil {
		return nil, err
	}
	vp := &DeleteVerbPhrase{
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

var delete_method_template = template.Must(
	template.New("delete_method_template").Parse(`
// {{.MethodName}} is part of the {{.InterfaceName}} interface.  defimpl verb {{.Verb.Tag}}.
func (x *{{.StructName}}) {{.MethodName}} (item {{.TypeString .SlotType.Elem}}) {
	i := -1
	for j, v := range x.{{.SlotName}} {
		if v == item {
			i = j
			break
		}
	}
	if i >= 0 {
		x.{{.SlotName}} = append(x.{{.SlotName}}[:i], x.{{.SlotName}}[i+1:]...)
	}
}
`))

// GlobalsTemplate is part of the VerbDefinition interface.
func (vd *Verb_Delete) GlobalsTemplate() *template.Template {
	return delete_method_template
}


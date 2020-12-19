package main

import "go/ast"
import "text/template"


type LengthVerbPhrase struct {
	slotVerbPhrase
}

var _ VerbPhrase = (*LengthVerbPhrase)(nil)
var _ SlotVerbPhrase = (*LengthVerbPhrase)(nil)
var _ GlobalsTemplateParameter = (*LengthVerbPhrase)(nil)


type Verb_Length struct {
	slotVerbDefinition
}

var _ VerbDefinition = (*Verb_Length)(nil)

func init() {
	vd := &Verb_Length{}
	VerbDefinitions[vd.Tag()] = vd
}

// Verb is part of the VerbDefinition interface.
func (vd *Verb_Length) Tag() string { return "length" }

// Description is part of the VerbDefinition interface.
func (vd *Verb_Length) Description() string {
	return "returns the length of the specified slice valued field."
}

// NewVerbPhrase is part of the VerbDefinition interface.
func (vd *Verb_Length) NewVerbPhrase(ctx *context, idef *InterfaceDefinition, field *ast.Field, comment *ast.Comment) (VerbPhrase, error) {
	slot, err := parse_slot_verb_phrase(ctx, field, comment)
	if err != nil {
		return nil, err
	}
	// The method signature for the length verb won't have a slot
	// type.  We dont need one though sice the method template
	// doesn't refer to the slot type.
	vp := &LengthVerbPhrase{
		slotVerbPhrase {
			baseVerbPhrase: baseVerbPhrase {
				verb: vd,
				idef: idef,
				field: field,
			},
			slot_name: slot,
			slot_type: nil,
		},
	}
	if err := addSlotSpec(idef, vp); err != nil {
		return nil, err
	}
	return vp, nil
}

var length_method_template = template.Must(
	template.New("length_method_template").Parse(`
// {{.MethodName}} is part of the {{.InterfaceName}} interface.
func (x *{{.StructName}}) {{.MethodName}}() int {
	return len(x.{{.SlotName}})
}
`))

// GlobalsTemplate is part of the VerbDefinition interface.
func (vd *Verb_Length) GlobalsTemplate() *template.Template {
	return length_method_template
}


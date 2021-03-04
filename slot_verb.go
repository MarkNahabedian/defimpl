package main

import "fmt"
import "reflect"
import "strings"
import "go/ast"
import "go/types"


// SlotVerbPhrase is an interface that should be implemented by all
// VerbPhrases that concern a slot of the impl struct of the
// VerbPhrase's InterfaceDefinition.
type SlotVerbPhrase interface {
	VerbPhrase
	SlotName() string 
	// Type returns the type of the slot.  For collection valued
	// slots (e.g. for the append/iterate, ... verbs) , this is
	// the Slot type of the collection itself (e.g. a slice),
	// rather than the element type.
	SlotType() types.Type
	SlotSpec() *slotSpec
	SetSlotSpec(*slotSpec)
	TypeString(t types.Type) string
}

type slotVerbPhrase struct {
	baseVerbPhrase
	slot_name string
	slot_spec *slotSpec
	slot_type types.Type
}

func (svp *slotVerbPhrase) Description() string {
	panic("(*slotVerbPhrase).Description() called.")
}

func (svp *slotVerbPhrase) Tag() string {
	panic("(*slotVerbPhrase).Tag() called.")
}

func (svp *slotVerbPhrase) GlobalDefinitions() (string, error) {
	panic("(*slotVerbPhrase).GlobalDefinitions() called.")
}

func (svp *slotVerbPhrase) SlotName() string {
	return svp.slot_name
}

func (svp *slotVerbPhrase) SlotSpec() *slotSpec {
	return svp.slot_spec
}

func (svp *slotVerbPhrase) SetSlotSpec(spec *slotSpec) {
	if svp.slot_spec != nil {
		panic("slot_spec already set.")
	}
	svp.slot_spec = spec
}

func (svp *slotVerbPhrase) SlotType() types.Type {
	return svp.slot_type
}

func (svp *slotVerbPhrase) TypeString(t types.Type) string {
	return types.TypeString(t, svp.InterfaceDefinition().File.Qualifier)
}


type slotSpec struct {
	VerbPhrases []SlotVerbPhrase
	// emitted is set to true when the slot declaration is added
	// to the impl struct.  This is so that only one slot is
	// defined no matter how namy VerbPhrases concern that slot.
	emitted bool
}

func (spec *slotSpec) InterfaceDefinition() *InterfaceDefinition {
	return spec.VerbPhrases[0].InterfaceDefinition()
}

func (spec *slotSpec) SlotName() string {
	return spec.VerbPhrases[0].(SlotVerbPhrase).SlotName()
}

func (spec *slotSpec) SlotType() types.Type {
	return spec.VerbPhrases[0].(SlotVerbPhrase).SlotType()
}


// addSlotSpec searches the InterfaceDefinition for a slotSpec with
// the same slot name as that of svp, and, failiing to find one,
// creates it, updating svf with the slotSpec.
func addSlotSpec(idef *InterfaceDefinition, svp SlotVerbPhrase) error {
	// *** IS THIS THE RIGHT TEST FOR TYPE EQUIVALENCE?
	teq := func(t1, t2 types.Type) bool {
		return (types.AssignableTo(t1, t2) &&
			types.AssignableTo(t2, t1))
	}
	for _, vp := range idef.VerbPhrases {
		if svp1, ok := vp.(SlotVerbPhrase); ok {
			if svp.SlotName() == svp1.SlotName() {
				// Verbs like length might not be able
				// to determine, nor need a SlotType.
				if svp.SlotType() != nil && !teq(svp.SlotType(), svp1.SlotType()) {
					return fmt.Errorf("Types %s and %s don't match",
						svp.SlotType(), svp1.SlotType())
				}
				svp.SetSlotSpec(svp1.SlotSpec())
				svp.SlotSpec().VerbPhrases = append(svp.SlotSpec().VerbPhrases, svp)
				return nil
			}
		}
	}
	svp.SetSlotSpec(&slotSpec {
		VerbPhrases: []SlotVerbPhrase { svp },
		emitted: false,
	})
	return nil
}


type slotVerbDefinition struct {}

func (vd *slotVerbDefinition) StructBody(vp VerbPhrase) (string, error) {
	svp := vp.(SlotVerbPhrase)
	if svp.SlotSpec() == nil {
		return "", fmt.Errorf("%#v has no SlotSpec", vp)
	}
	if !svp.SlotSpec().emitted {
		svp.SlotSpec().emitted = true
		return fmt.Sprintf("\t%s %s\n", svp.SlotName(),
			svp.TypeString(svp.SlotType())), nil
	}
	return "", nil
}


// parse_slot_verb_phrase parses the defimpl comment for verbs that
// parse to aa SlotVerbPhrase.
func parse_slot_verb_phrase(ctx *context, field *ast.Field, comment *ast.Comment) (string, error) {
	val, ok := reflect.StructTag(comment.Text[2:]).Lookup("defimpl")
	if !ok {
		// Shouldn't happen.  To get here we should already have found a defimpl comment.
		panic("Can't construct VerbPhrase from a non-defimpl comment.")
	}
	split := strings.Split(val, " ")
	if len(split) != 2 {
		pos := ctx.fset.Position(comment.Slash)
		return "", fmt.Errorf("defimpl verb %q expects 1 parameter, a slot name: %s %q",
			split[0], pos, comment.Text)
	}
	if len(field.Names) != 1 {
		pos := ctx.fset.Position(comment.Slash)
		return "", fmt.Errorf("defimpl verb %q requires a field with only one name: %s",
			split[0], pos)
	}
	return split[1], nil
}


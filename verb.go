package main

import "bytes"
import "fmt"
import "os"
import "reflect"
import "strings"
import "go/ast"
import "text/template"


// GlobalsTemplateParameter codifies the parameter of the
// (VerbDefinition).GlobalsTemplate.
type GlobalsTemplateParameter interface {
	MethodName() string
	InterfaceName() string
	StructName() string
	SlotName() string
}

// VerbPhrase represents a single verb and its parameters from a Field
// in an interface definition.
type VerbPhrase interface {
	// Verb returns the VerbDefinition that created this VerbPhrase.
	Verb() VerbDefinition
	// Tag returns the string that identifies a Verb in a tag
	// comment.
	Tag() string
	// MethodName returns the function name of the method
	// associated with this verb phrase.
	MethodName() string
	// InterfaceDefinition returns the InterfaceDefinition which
	// represents the interface type declaration that this
	// VerbPhrase's Field appears in.
	InterfaceDefinition() *InterfaceDefinition
	// Field returns the ast.Field whose tag comment the verb
	// phrase is derived from.
	Field() *ast.Field

	// GlobalDefinitions returns the global code that should be
	// included to support the VerbPhrase.
	// GlobalDefinitions() (string, error)

	// The following methods are a convenience for implementing
	// GlobalsTemplate templates.  Such templates are executed with
	// either a VerbPhrase or a GlobalsTemplateParameter (see
	// CheckSignatures).

	// InterfaceName returns the InterfaceName from the InterfaceDefinition.
	InterfaceName() string
	// StructName returns the StructName from the InterfaceDefinition
	StructName() string
}

func GetVerbPhrase(ctx *context, idef *InterfaceDefinition, method *ast.Field) {
	if method.Comment == nil {
		return
	}
	for _, c := range method.Comment.List {
		val, ok := reflect.StructTag(c.Text[2:]).Lookup("defimpl")
		if !ok {
			continue
		}
		split := strings.Split(val, " ")
		if len(split) < 1 {
			continue
		}
		// constructor, ok := Verbs[split[0]]
		vd, ok := VerbDefinitions[split[0]]
		if !ok {
			fmt.Fprintf(os.Stderr, "defimpl: Unknown verb %q in defimpl comment %s: %q\n",
				split[0], ctx.fset.Position(c.Slash), c.Text)
			continue
		}
		// vp, err := constructor(ctx, idef, method, c)
		vp, err := vd.NewVerbPhrase(ctx, idef, method, c)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		} else {
			if vp != nil {
				idef.VerbPhrases = append(idef.VerbPhrases, vp)
			}
		}
	}
}

// GlobalDefinitions returns the definition of the method for this verb
// based on the Verb.GlobalsTemplate.
func GlobalDefinitions(vp VerbPhrase) (string, error) {
	tmpl := vp.Verb().GlobalsTemplate()
	if tmpl == nil {
		return "", nil
	}
	w := &bytes.Buffer{}
	if err := tmpl.Execute(w, vp); err != nil {
		return "", err
	}
	return w.String(), nil
}


type VerbDefinition interface {
	Tag() string
	Description() string
	NewVerbPhrase(*context, *InterfaceDefinition, *ast.Field, *ast.Comment) (VerbPhrase, error)
	GlobalsTemplate() *template.Template
	StructBody(VerbPhrase) (string, error)
}

var VerbDefinitions map[string]VerbDefinition = map[string]VerbDefinition{}



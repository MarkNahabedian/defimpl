package main

import "go/ast"
import "text/template"


type DiscriminateVerbPhrase struct {
	baseVerbPhrase
}

var _ VerbPhrase = (*DiscriminateVerbPhrase)(nil)
// var _ GlobalsTemplateParameter = (*DiscriminateVerbPhrase)(nil)

func (vp *DiscriminateVerbPhrase)StructBody() (string, error) {
	return "", nil
}


type Verb_Discriminate struct {}

var _ VerbDefinition = (*Verb_Discriminate)(nil)

func init() {
	vd := &Verb_Discriminate{}
	VerbDefinitions[vd.Tag()] = vd
}

// Verb is part of the VerbDefinition interface.
func (vd *Verb_Discriminate) Tag() string { return "discriminate" }

// Description is part of the VerbDefinition interface.
func (vd *Verb_Discriminate) Description() string {
	return "the empty method that distinguishes implementors of this interface from those that would otherwise have the same method set."
}

// NewVerbPhrase is part of the VerbDefinition interface.
func (vd *Verb_Discriminate) NewVerbPhrase(ctx *context, idef *InterfaceDefinition, field *ast.Field, comment *ast.Comment) (VerbPhrase, error) {
	_, err := CheckSignatures(ctx, vd, idef.Package(), field, vd.GlobalsTemplate())
	if err != nil {
		return nil, err
	}
	vp := &DiscriminateVerbPhrase{
		baseVerbPhrase: baseVerbPhrase {
			verb: vd,
			idef: idef,
			field: field,
		},
	}
	return vp, nil
}

var discriminate_method_template = template.Must(
		template.New("discriminate_method_template").Parse(`
// {{.MethodName}} is part of the {{.InterfaceName}} interface.  defimpl verb {{.Verb.Tag}}.
func (x *{{.StructName}}) {{.MethodName}}() {}
`))

// GlobalsTemplate is part of the VerbDefinition interface.
func (vd *Verb_Discriminate) GlobalsTemplate() *template.Template {
	return discriminate_method_template
}

func (vd *Verb_Discriminate) StructBody(VerbPhrase) (string, error) {
	return "", nil
}


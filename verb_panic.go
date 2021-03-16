package main

import "go/ast"
import "text/template"
import "defimpl/util"


type PanicVerbPhrase struct {
	baseVerbPhrase
	MethodParameters string
	MethodResults string
}

var _ VerbPhrase = (*PanicVerbPhrase)(nil)

func (vp *PanicVerbPhrase) StructBody() (string, error) {
	return "", nil
}


type Verb_Panic struct {}

var _ VerbDefinition = (*Verb_Panic)(nil)

func init() {
	vd := &Verb_Panic{}
	VerbDefinitions[vd.Tag()] = vd
}

// Tag is part of the VerbDefinition interface.
func (vd *Verb_Panic) Tag() string { return "panic" }

// Description is part of the VerbDefinition interface.
func (vd *Verb_Panic) Description() string {
	return "the method will panic if called, for when an implementation only needs to partially implement an interface."
}

// NewVerbPhrase is part of the VerbDefinition interface.
func (vd *Verb_Panic) NewVerbPhrase(ctx *context, idef *InterfaceDefinition, field *ast.Field, comment *ast.Comment) (VerbPhrase, error) {
	_, err, scratchpad := CheckSignatures(ctx, vd, idef.Package(), field, vd.GlobalsTemplate())
	if err != nil {
		return nil, err
	}
	q := util.TypeStringQualifier(idef.File.AstFile)
	params := ""
	if p, ok := scratchpad["__PARAMETERS"].(*ast.FieldList); ok {
		params, _ = util.FieldListString(p, ctx.info, q, false, false)
	}
	results := ""
	if r, ok := scratchpad["__RESULTS"].(*ast.FieldList); ok {
		results, _ = util.FieldListString(r, ctx.info, q, false, true)
	}
	vp := &PanicVerbPhrase{
		baseVerbPhrase: baseVerbPhrase {
			verb: vd,
			idef: idef,
			field: field,
		},
		MethodParameters: params,
		MethodResults: results,
	}
	return vp, nil
}

var panic_method_template = template.Must(
		template.New("panic_method_template").Parse(`
// {{.MethodName}} is part of the {{.InterfaceName}} interface.  defimpl verb {{.Verb.Tag}}.
func (x *{{.StructName}}) {{.MethodName}}({{.MethodParameters}}) {{.MethodResults}} {
	panic("(*{{.StructName}}).{{.MethodName}} was called")
}
`))

// GlobalsTemplate is part of the VerbDefinition interface.
func (vd *Verb_Panic) GlobalsTemplate() *template.Template {
	return panic_method_template
}

func (vd *Verb_Panic) StructBody(VerbPhrase) (string, error) {
	return "", nil
}


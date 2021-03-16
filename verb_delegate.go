// The delegate verb is meant as an aid to developing an an interface
// impleementation that can delegate some interface methods to some
// data member.
//
// The parameter of the delegate verb is the sewcond half of a
// selector expression that will be combined with the receiver to
// identify where to delegate the operation to.

package main

import "fmt"
import "reflect"
import "strings"
import "go/ast"
import "text/template"
import "defimpl/util"


type DelegateVerbPhrase struct {
	baseVerbPhrase
	MethodParameters string
	ParameterNames string
	MethodResults string
	DelegateTo string
}

var _ VerbPhrase = (*DelegateVerbPhrase)(nil)

func (vp *DelegateVerbPhrase) StructBody() (string, error) {
	return "", nil
}


type Verb_Delegate struct {}

var _ VerbDefinition = (*Verb_Delegate)(nil)

func init() {
	vd := &Verb_Delegate{}
	VerbDefinitions[vd.Tag()] = vd
}

// Tag is part of the VerbDefinition interface.
func (vd *Verb_Delegate) Tag() string { return "delegate" }

// Description is part of the VerbDefinition interface.
func (vd *Verb_Delegate) Description() string {
	return "the method will delegate to another object."
}

// NewVerbPhrase is part of the VerbDefinition interface.
func (vd *Verb_Delegate) NewVerbPhrase(ctx *context, idef *InterfaceDefinition, field *ast.Field, comment *ast.Comment) (VerbPhrase, error) {
	_, err, scratchpad := CheckSignatures(ctx, vd, idef.Package(), field, vd.GlobalsTemplate())
	if err != nil {
		return nil, err
	}
	delegate_to, err := parse_DelegateVerbPhrase(ctx, field, comment)
	if err != nil {
		return nil, err
	}
	q := util.TypeStringQualifier(idef.File.AstFile)
	// For method parameters we need two forms:
	//   * the parameter list declaration
	///  * the parameter names to pass to the delegate.
	params := ""
	actual := ""
	if p, ok := scratchpad["__PARAMETERS"].(*ast.FieldList); ok {
		params, actual = util.FieldListString(p, ctx.info, q, true, false)
	}
	results := ""
	if r, ok := scratchpad["__RESULTS"].(*ast.FieldList); ok {
		results, _ = util.FieldListString(r, ctx.info, q, false, true)
	}
	vp := &DelegateVerbPhrase{
		baseVerbPhrase: baseVerbPhrase {
			verb: vd,
			idef: idef,
			field: field,
		},
		MethodParameters: params,
		ParameterNames: actual,
		MethodResults: results,
		DelegateTo: delegate_to,
	}
	return vp, nil
}

func parse_DelegateVerbPhrase(ctx *context, field *ast.Field, comment *ast.Comment) (string, error) {
	val, ok := reflect.StructTag(comment.Text[2:]).Lookup("defimpl")
	if !ok {
		// Shouldn't happen.  To get here we should already have found a defimpl comment.
		panic("Can't construct VerbPhrase from a non-defimpl comment.")
	}
	// Everything after the verb itself is captured as the
	// delation target, in case the string representation of that
	// target contains a space.
	split := strings.SplitN(val, " ", 2)
	if len(split) != 2 {
		pos := ctx.fset.Position(comment.Slash)
		return "", fmt.Errorf("defimpl verb %q expects 1 parameter, the right hand side of a selector expression: %s %q",
			split[0], pos, comment.Text)
	}
	return split[1], nil
}

var delegate_method_template = template.Must(
		template.New("delegate_method_template").Parse(`
// {{.MethodName}} is part of the {{.InterfaceName}} interface.  defimpl verb {{.Verb.Tag}}.
func (x *{{.StructName}}) {{.MethodName}}({{.MethodParameters}}) ({{.MethodResults}}) {
	return x.{{.DelegateTo}}.{{.MethodName}}({{.ParameterNames}})
}
`))

// GlobalsTemplate is part of the VerbDefinition interface.
func (vd *Verb_Delegate) GlobalsTemplate() *template.Template {
	return delegate_method_template
}

func (vd *Verb_Delegate) StructBody(VerbPhrase) (string, error) {
	return "", nil
}



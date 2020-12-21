package main

import "bytes"
import "fmt"
import "reflect"
import "strings"
import "go/ast"
import "go/parser"
import "go/types"
import "text/template"
import "defimpl/util"


type EmbedVerbPhrase struct {
	baseVerbPhrase
	ImplStruct string
	EmbeddedInterface string
}

var _ VerbPhrase = (*EmbedVerbPhrase)(nil)
// var _ GlobalsTemplateParameter = (*EmbedVerbPhrase)(nil)

func (vp *EmbedVerbPhrase)StructBody() (string, error) {
	return "", nil
}


type Verb_Embed struct {}

var _ VerbDefinition = (*Verb_Embed)(nil)

func init() {
	vd := &Verb_Embed{}
	VerbDefinitions[vd.Tag()] = vd
}

// Verb is part of the VerbDefinition interface.
func (vd *Verb_Embed) Tag() string { return "embed" }

// Description is part of the VerbDefinition interface.
func (vd *Verb_Embed) Description() string {
	return "Specifies a concrete type to embed to implement an interface."
}

// NewVerbPhrase is part of the VerbDefinition interface.
func (vd *Verb_Embed) NewVerbPhrase(ctx *context, idef *InterfaceDefinition, field *ast.Field, comment *ast.Comment) (VerbPhrase, error) {
	// We expect the method signature to have an interface type
	// but no name.
	//
	// The verb comment might have an optional parameter that is
	// the type to embed.  It might come from a package that is
	// not otherwise used by the input file.  In that case we
	// expect that package to be imported using the _ identifier.

	// For us to have gotten here, this must have worked
	// previously:
	structTag, _:= reflect.StructTag(comment.Text[2:]).Lookup("defimpl")

	impl := ""
	split := strings.Split(structTag, " ")
	if len(split) > 1 {
		v := ParseEmbedImpl(split[1])
		if v.err != nil {
			return nil, v.err
		}
		impl = v.ImplName()
	}
	embedded_package := ""
	embedded_name := ""
	switch e := field.Type.(type) {
		case *ast.Ident:
			embedded_package = ""
			embedded_name = e.Name
		case *ast.SelectorExpr:
			embedded_package = types.ExprString(e.X)
			embedded_name = e.Sel.Name

		default:
			panic(fmt.Sprintf("Unsupported EXpr type %T", field.Type))
		}
	// If no struct specified then assume that the embedded
	// interface has a defimpl generated struct:
	if impl == "" {
		impl = util.ImplName(embedded_package, embedded_name)
	}
	vp := &EmbedVerbPhrase{
		baseVerbPhrase: baseVerbPhrase {
			verb: vd,
			idef: idef,
			field: field,
		},
		ImplStruct: impl,
		EmbeddedInterface: "",
	}
	if embedded_package == "" {
		vp.EmbeddedInterface = embedded_name
	} else {
		vp.EmbeddedInterface = embedded_package + "." + embedded_name
	}
	return vp, nil
}


var embed_method_template = template.Must(
	template.New("embed_method_template").Parse(`
var _ {{.EmbeddedInterface}} = (*{{.StructName}})(nil)  // defimpl verb {{.Verb.Tag}}.
`))

// GlobalsTemplate is part of the VerbDefinition interface.
func (vd *Verb_Embed) GlobalsTemplate() *template.Template {
	return embed_method_template
}

func (vd *Verb_Embed) StructBody(vp VerbPhrase) (string, error) {
	return vp.(*EmbedVerbPhrase).ImplStruct, nil
}

type parse_embed_visitor struct {
	pointer bool
	identifiers []string
	err error
}

func (v *parse_embed_visitor) Path() string {
	return strings.Join(v.identifiers[0:len(v.identifiers) - 1], "/")
}

func (v *parse_embed_visitor) Name() string {
	return v.identifiers[len(v.identifiers) - 1]
}

func (v *parse_embed_visitor) ImplName() string {
	w := bytes.NewBufferString("")
	if v.pointer {
		w.Write([]byte("*"))
	}
	length := len(v.identifiers)
	if length > 1 {
		w.Write([]byte(v.identifiers[length - 2]))
		w.Write([]byte("."))
	}
	w.Write([]byte(v.identifiers[length - 1]))
	return w.String()
}

func (v *parse_embed_visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}
	switch n1 := n.(type) {
	case *ast.StarExpr:
		v.pointer = true
	case *ast.Ident:
		v.identifiers = append(v.identifiers, n1.Name)
	}
	return v
}

func ParseEmbedImpl(s string) *parse_embed_visitor {
	visitor := &parse_embed_visitor {
		pointer: false,
		identifiers: []string{},
		err: nil,
	}
	expr, err := parser.ParseExpr(s)
	if err != nil {
		visitor.err = err
		return visitor
	}
	ast.Walk(visitor, expr)
	return visitor
}


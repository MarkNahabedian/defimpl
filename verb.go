package main

import "bytes"
import "defimpl/util"
import "fmt"
import "go/ast"
import "text/template"

type VerbTemplateParameter struct {
	Verb *VerbDefinition
	InterfaceName string
	StructName string
	MethodName string
	SlotName string
	Type ast.Expr  // types.Type
}

func (v *VerbTemplateParameter) RunTemplate() string {
	writer := bytes.NewBufferString("")
	err := v.Verb.Template.Execute(writer, v)
	if err != nil {
		panic(err)
	}
	return writer.String()
}

type VerbDefinition struct {
	Verb string
	Description string
	// Assimilate assimilates the method into the slotSpec if appropriate.
	Assimilate func(*context, *VerbDefinition, *slotSpec, *InterfaceDefinition, *ast.Field) error
	// Template will generate the code associated with this Verb.
	// The template will be passed a VerbTemplateParameter as its parameter.
	Template *template.Template
}

var VerbDefinitions map[string]*VerbDefinition = map[string]*VerbDefinition{}

func LookupVerb(verb string) *VerbDefinition {
	vd, ok := VerbDefinitions[verb]
	if ok {
		return vd
	}
	return nil
}

// checkSignature returns an error if fd has the wrong number of parameters
// or return values.
func checkSignature(fd *ast.FuncType, paramCount, resultCount int) error {
	if len(util.FieldListSlice(fd.Params)) != paramCount ||
		len(util.FieldListSlice(fd.Results)) != resultCount {
		return fmt.Errorf("alleged method has inappropriate signature\n")
	}
	return nil
}

func SliceOfType(typ ast.Expr) ast.Expr {
	return &ast.ArrayType{ Elt: typ }
}

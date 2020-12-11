package main

import "go/ast"


type baseVerbPhrase struct {
	verb VerbDefinition
	idef *InterfaceDefinition
	field *ast.Field
}

// Verb is part of the VerbPhrase interface.
func (vp *baseVerbPhrase) Tag() string {
	return vp.verb.Tag()
}

// Verb is part of the VerbPhrase interface.
func (vp *baseVerbPhrase) Verb() VerbDefinition {
	return vp.verb
}

// InterfaceDefinition is part of the VerbPhrase interface.
func (vp *baseVerbPhrase) InterfaceDefinition() *InterfaceDefinition {
	return vp.idef
}

// Field is part of the VerbPhrase interface.
func (vp *baseVerbPhrase) Field() *ast.Field {
	return vp.field
}

func (vp *baseVerbPhrase) MethodName() string {
	return vp.field.Names[0].Name
}

func (vp *baseVerbPhrase) InterfaceName() string {
	return vp.InterfaceDefinition().InterfaceName
}

func (vp *baseVerbPhrase) StructName() string {
	return vp.InterfaceDefinition().StructName()
}


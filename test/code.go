// This is a test program that uses defimpl to generate struct
// definitions that provide default implementations of interfaces.
// The program verifies that the generated code functions properly.
package test

import "reflect"
import tmpl "text/template"

//go:generate defimpl

type Thing interface {
	// name
	Name() string   // defimpl:"read name"
	SetName(string) // defimpl:"set name"
	// related
	AddRelated(...Thing)        // defimpl:"append related"
	GetRelated(int) Thing       // defimpl:"index related"
	CountRelated() int          // defimpl:"length related"
	DoRelated(func(Thing) bool) // defimpl:"iterate related"

	// These are added to test that the proper packages are
	// imported in the output file.

	// mytype
	MyType() reflect.Type // defimpl:"read mytype"
	SetType(reflect.Type) // defimpl:"set mytype"
	// template
	Template() *tmpl.Template // defimpl:"read template"
}

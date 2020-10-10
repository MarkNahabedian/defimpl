// This is a test program that uses defimpl to generate struct
// definitions that provide default implementations of interfaces.
// The program verifies that the generated code functions properly.
package test

import "reflect"
import tmpl "text/template"

//go:generate defimpl

type Thing interface {
	// name
	Thing()          // defimpl:"discriminate"
	Name() string    // defimpl:"read name"
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

// Base1 (ABSTRACT)
type Base1 interface {
	Id() int  // defimpl:"read id"
}

// Base2 (ABSTRACT)
type Base2 interface {
	Name() string  // defimpl:"read name"
}

// Sub1 (ABSTRACT)
type Sub1 interface {
	Base1
	Color() string  // defimpl:"read color"
}

type Gazong interface {
	Base2
	Sub1
	Smug() bool     // defimpl:"read smug"
}

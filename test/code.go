// This is a test program that uses defimpl to generate struct
// definitions that provide default implementations of interfaces.
// The program verifies that the generated code functions properly.
package test

//go:generate defimpl

type Thing interface {
	Name() string               // defimpl:"read name"
	SetName(string)             // defimpl:"set name"
	AddRelated(...Thing)        // defimpl:"append related"
	GetRelated(int) Thing       // defimpl:"index related"
	CountRelated() int          // defimpl:"length related"
	DoRelated(func(Thing) bool) // defimpl:"iterate related"
}

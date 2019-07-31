// A simple example of using defimpl.
package test

type Tower interface {
	Height() float   // defimpl:"read height"
}

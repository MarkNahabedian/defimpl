// A simple example of using defimpl.
package test

type Tower interface {
	Height() float32   // defimpl:"read height"
}

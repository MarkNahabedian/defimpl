package test

import "testing"

func NewThing() Thing {
	return &ThingImpl{}
}

func TestReadSet(t *testing.T) {
	thing1 := NewThing()
	if want, got := "", thing1.Name(); want != got {
		t.Errorf("Reading name of unnamed thing: got: %v, want: %v", got, want)
	}
	thing1.SetName("thing1")
	if want, got := "thing1", thing1.Name(); want != got {
		t.Errorf("Reading name of thing: got %v, want %v", got, want)
	}
}

/*
func TestAppendIterate(t *testing.T) {
	thing1 := NewThing()
	thing1.SetName("thing1")
	thing2 := NewThing()
	thing2.SetName("thing2")
	thing3 := NewThing()
	thing3.SetName("thing3")
	thing1.AddRelated(thing2)
	thing1.AddRelated(thing3)
	expect := []Thing{ thing2, thing3 }
	expect_index := 0
	thing1.DoRelated(func(thing Thing) bool {
		if want := expect[expect_index]; thing != want {
			t.Errorf("%d: got %#v, want %#v", expect_index, thing, want)
		}
		expect_index += 1
		return true
	})
}
*/

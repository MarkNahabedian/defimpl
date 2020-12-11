package test

import "testing"

func NewThing() Thing {
	return Thing(&ThingImpl{})
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

func TestSliceValued(t *testing.T) {
	thing1 := NewThing()
	thing1.SetName("thing1")
	thing2 := NewThing()
	thing2.SetName("thing2")
	thing3 := NewThing()
	thing3.SetName("thing3")
	thing1.AddRelated(thing2)
	thing1.AddRelated(thing3)
	if want, got := 2, thing1.CountRelated(); got != want {
		t.Errorf("Wrong slice length, got %d, want %d", got, want)
	}
	if want, got := thing3, thing1.GetRelated(1); got != want {
		t.Errorf("Wrong element at index 1: got %#v, want %#v", got, want)
	}
	test_iterate := func(expect []Thing) {
		expect_index := 0
		thing1.DoRelated(func(thing Thing) bool {
			if want := expect[expect_index]; thing != want {
				t.Errorf("%d: got %#v, want %#v", expect_index, thing, want)
			}
			expect_index += 1
			return true
		})
	}
	test_iterate([]Thing{thing2, thing3})
	thing1.RemoveRelated(thing2)
	test_iterate([]Thing{thing3})
}

func TestInheritance(t *testing.T) {
}

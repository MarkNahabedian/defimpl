package test

import "reflect"
import "testing"
import "defimpl/runtime"

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

func TestRuntime(t *testing.T) {
	thing := &ThingImpl{}
	thing.SetName("foo")
	ty := reflect.TypeOf(thing)
	if ty.Kind() != reflect.Ptr || ty.Elem().Kind() != reflect.Struct {
		t.Errorf("ThingImpl is %v, not pointer to struct", ty)
	}
	if want, got := ty, runtime.ImplFor(ty); got != want {
		t.Errorf("ImplFor of Impl type failed: want %v, got %v", want, got)
	}
	i := runtime.InterfaceFor(ty)
	if want, got := reflect.Interface, i.Kind(); want != got {
		t.Errorf("InterfaceFor of Impl type failed: want %v, got %v", want, got)
	}
	iimpl := runtime.ImplFor(i)
	if want, got := reflect.Ptr, iimpl.Kind(); want != got {
		t.Errorf("ImplFor of interface type failed: want %v, got %v", want, got)
	}
	if want, got := reflect.Struct, iimpl.Elem().Kind(); want != got {
		t.Errorf("ImplFor of interface type failed: want %v, got %v", want, got)
	}
}


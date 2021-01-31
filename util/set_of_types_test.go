package util

import "reflect"
import "testing"

type struct1 struct {
	bool
	uint
}

func Test_SetOfTypes(t *testing.T) {
	types := MemberTypes(reflect.TypeOf(
		struct {
			int
			string
			struct1
		}{}))
	if typ := reflect.TypeOf(true); types.Contains(typ) {
		t.Errorf("%v contains %v, but it shouldn't.", types, typ)
	}
	if typ := reflect.TypeOf(struct1{}); !types.Contains(typ) {
		t.Errorf("%v should contain %v, but it does't.", types, typ)
	}

}


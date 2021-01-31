// set_of_types provides utilities for dealing with sets of Go types.
package util

import "reflect"


// SetOfTypes implements a set of reflect.Type.
type SetOfTypes []reflect.Type


// Contains returns true if t appears in the set.
func (st SetOfTypes) Contains(t reflect.Type) bool {
	for _, t1 := range st {
		if t1 == t {
			return true
		}
	}
	return false
}


// MemberTypes returns a SetOfTypes.  The parameter should be a struct
// type or a pointer to such.  The returned SetOfTypes will consist of
// each type that appears as a member of the struct.
func MemberTypes(t reflect.Type) SetOfTypes {
	result := SetOfTypes{}
	include := func (t reflect.Type) {
		for _, present := range result {
			if present == t {
				return
			}
		}
		result = append(result, t)
	}
	switch t.Kind() {
	case reflect.Ptr:
		return MemberTypes(t.Elem())
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			include(t.Field(i).Type)
		}
	}
	return result
}


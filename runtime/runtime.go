// Package runtime encompases the runtime data created by the defimpl
// utility.

package runtime

import "fmt"
import "reflect"


// interfaceToImpl maps from the interfaces that defimpl has defined
// implementation structs for to the corresponding struct pointer types.
var interfaceToImpl = map[reflect.Type]reflect.Type{}

//iImplToInterface maps from an implementation type defined by the
// defimpl utility to the interface type it is defined for.
var implToInterface = map[reflect.Type]reflect.Type{}

// InterfaceToImpl returns the implementation type (as defined by
// defimpl) for the specified interface type.
func InterfaceToImpl(inter reflect.Type) reflect.Type {
	return interfaceToImpl[inter]
}

// ImplToInterface returns the interface type that the type impl was
// defined for.
func ImplToInterface(impl reflect.Type) reflect.Type {
	return implToInterface[impl]
}

// InterfaceFor returns the iinterface type for the specified type,
// assuming that they are under the perview of defimpl.
func InterfaceFor(t reflect.Type) (reflect.Type, error) {
	switch t.Kind() {
	case reflect.Interface:
		return t, nil
	case reflect.Ptr:
		if t.Elem().Kind() == reflect.Struct {
			if i, ok := implToInterface[t]; ok {
				return i, nil
			} else {
				return nil, fmt.Errorf("%s not found in implToInterface", t)
			}
		}
		return nil, fmt.Errorf("%s is pointer, but not to struct", t)
	case reflect.Struct:
		return ImplToInterface(reflect.PtrTo(t)), nil
	}
	return nil, fmt.Errorf("%s is neither interface nor pointer to struct", t)
}

// ImplFor returns the iinterface type for the specified type,
// assuming that they are under the perview of defimpl.
func ImplFor(t reflect.Type) (reflect.Type, error) {
	switch t.Kind() {
	case reflect.Ptr:
		if t.Elem().Kind() == reflect.Struct {
			return t, nil
		}
		return nil, fmt.Errorf("%s is a pointer, but not to struct", t)
	case reflect.Interface:
		return InterfaceToImpl(t), nil
	}
	return nil, fmt.Errorf("ImplFor called on  %s but expected interface or pointer to struct type", t)
}


// Register associates the specified interface and implementation
// types so that they can be found at run time.
//
// Regiuster should only be called from init functions defined in the
// code generated by the defimpl utilitty.
//
// Register should only be called from code generated by defimpl.
// Panicing is appropriate for unexpected situations.
func Register(inter reflect.Type, impl reflect.Type) {
	if inter == nil {
		panic(fmt.Sprintf("defimpl/runtime.Register(%v, %v): interface is nil", inter, impl))
	}
	if impl == nil {
		panic(fmt.Sprintf("defimpl/runtime.Register(%v, %v): impl is nil", inter, impl))
	}
	if inter.Kind() != reflect.Interface {
		panic(fmt.Sprintf("%v should be an interface, not %s",
			inter, inter.Kind().String()))
	}
	if impl.Kind() == reflect.Ptr {
		if  impl.Elem().Kind() != reflect.Struct {
			panic(fmt.Sprintf("%v should be pointer to struct, not %s",
				impl, impl.Kind().String()))
		}
	} else {
		panic(fmt.Sprintf("%v should be pointer (to struct), not %s",
			impl, impl.Kind().String()))
	}
	interfaceToImpl[inter] = impl
	implToInterface[impl] = inter
}

func Dump() {
	fmt.Println("interfaceToImpl:")
	for k, v := range interfaceToImpl {
		fmt.Println("\t", k, "\t", v)
	}
	fmt.Println("implToInterface:")
	for k, v := range implToInterface {
		fmt.Println("\t", k, "\t", v)
	}
}


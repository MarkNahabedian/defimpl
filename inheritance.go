// Process interface inheritance.
package main

import "fmt"
import "os"

// DoInheritance fills in the AllInherited field of
// InterfaceDefinition while a context is readily available.
func (ctx *context) DoInheritance() {
	for _, f := range ctx.files {
		f.DoInheritance(ctx)
	}
}

func (f *File) DoInheritance(ctx *context) {
	for _, i := range f.Interfaces {
		i.DoInheritance(ctx)
	}
}

func (idef *InterfaceDefinition) DoInheritance(ctx *context) {
	_ = idef.GetInherited(ctx)
}

func (idef *InterfaceDefinition) GetInherited(ctx *context) []*InterfaceDefinition {
	var gi func(*InterfaceDefinition, []*InterfaceDefinition) []*InterfaceDefinition
	gi = func(idef *InterfaceDefinition, circular []*InterfaceDefinition) []*InterfaceDefinition {
		if idef.AllInherited != nil {
			return idef.AllInherited
		}
		for _, c := range circular {
			if idef == c {
				fmt.Fprintf(os.Stderr, "defimpl: Circular interface definitions: %s within %s",
					idef.QualifiedName(), c.QualifiedName())
				return []*InterfaceDefinition{}
			}
		}
		circular = append(circular, idef)
		all_inherited := []*InterfaceDefinition{}
		adjoin := func(id *InterfaceDefinition) {
			for _, i := range all_inherited {
				if i == id {
					return
				}
			}
			all_inherited = append(all_inherited, id)
		}
		for _, inherited := range idef.Inherited {
			ih := ctx.IDLookup(inherited)
			if ih == nil {
				fmt.Fprintf(os.Stderr, "defimpl: For interface %s: Can't find inherited interface %s.\n",
					idef.QualifiedName(), inherited)
				continue
			}
			adjoin(ih)
			for _, inherited := range gi(ih, circular) {
				adjoin(inherited)
			}
		}
		idef.AllInherited = all_inherited
		return idef.AllInherited
	}
	return gi(idef, []*InterfaceDefinition{})
}

func (idef *InterfaceDefinition) InheritedVerbs() /* []*VerbTemplateParameter */ {
	panic("NYI")
	/*
	result := []*VerbTemplateParameter{}
	for _, inherited := range idef.AllInherited {
		if !inherited.IsAbstract {
			continue
		}
		for _, sspec := range inherited.SlotSpecs {
			for _, v := range sspec.Verbs {
				copied := *v
				copied.StructName = idef.StructName()
				result = append(result, &copied)
			}
		}
	}
	*/
	// return result
}

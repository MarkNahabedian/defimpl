package main

import "defimpl/util"
import "go/ast"
import "go/types"
import "text/template"

func init() {
	vd := &VerbDefinition{Verb: "append"}
	vd.Description = "appends the specified value to the field."
	vd.Assimilate = func(ctx *context, vd *VerbDefinition, spec *slotSpec, id *InterfaceDefinition, m *ast.Field) error {
		ftype, ok := m.Type.(*ast.FuncType)
		if !ok {
			return nil
		}
		if err := checkSignature(ftype, 1, 0); err != nil {
			return err
		}
		params := util.FieldListSlice(ftype.Params)
		e, ok := params[0].Type.(*ast.Ellipsis)
		if !ok {
			return nil
		}
		spec.CheckType(SliceOfType(e.Elt))
		spec.Verbs = append(spec.Verbs, &VerbTemplateParameter{
			Verb:          vd,
			InterfaceName: id.InterfaceName,
			StructName:    id.StructName(),
			MethodName:    m.Names[0].Name,
			SlotName:      spec.Name,
			Type:          spec.Type,
		})
		return nil
	}
	vd.Template = template.Must(template.New(vd.Verb).Funcs(map[string]interface{}{
		"ExprString": types.ExprString,
	}).Parse(`
		func (x *{{.StructName}}) {{.MethodName}} (v ...{{ExprString .Type.Elt}}) {
			x.{{.SlotName}} = append(x.{{.SlotName}}, v...)
		}
	`))
	VerbDefinitions[vd.Verb] = vd
}

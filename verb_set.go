package main

import "defimpl/util"
import "go/ast"
import "text/template"

func init() {
	vd := &VerbDefinition{ Verb: "set" }
	vd.Description = "sets the value of the field to that provided."
	vd.Assimilate = func(ctx *context, vd *VerbDefinition, spec *slotSpec, id *InterfaceDefinition, m *ast.Field) error {
		ftype, ok := m.Type.(*ast.FuncType)
		if !ok {
			return nil
		}
		if err := checkSignature(ftype, 1, 0); err != nil {
			return err
		}
		params := util.FieldListSlice(ftype.Params)
		spec.CheckType(params[0].Type)
		spec.Verbs = append(spec.Verbs, &VerbTemplateParameter{
			Verb: vd,
			InterfaceName: id.InterfaceName,
			StructName: id.StructName(),
			MethodName: m.Names[0].Name,
			SlotName: spec.Name,
			Type: spec.Type,
		})
		return nil
	}
	vd.Template = template.Must(template.New(vd.Verb).Parse(`
		func (x *{{.StructName}}) {{.MethodName}}(v {{.Type.String}}) {
			x.{{.SlotName}} = v
		}
	`))
	VerbDefinitions[vd.Verb] = vd
}

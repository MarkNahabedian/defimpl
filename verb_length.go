package main

import "defimpl/util"
import "go/ast"
import "go/types"
import "text/template"

func init() {
	vd := &VerbDefinition{Verb: "length"}
	vd.Description = "returns the length of the specified slice valued field."
	vd.Assimilate = func(ctx *context, vd *VerbDefinition, spec *slotSpec, id *InterfaceDefinition, m *ast.Field) error {
		ftype, ok := m.Type.(*ast.FuncType)
		if !ok {
			return nil
		}
		if err := checkSignature(ftype, 0, 1); err != nil {
			return err
		}
		results := util.FieldListSlice(ftype.Results)
		spec.CheckType(SliceOfType(results[0].Type))
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
		func (x *{{.StructName}}) {{.MethodName}} () int {
			return len(x.{{.SlotName}})
		}
	`))
	VerbDefinitions[vd.Verb] = vd
}

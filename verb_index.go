package main

import "defimpl/util"
import "fmt"
import "go/ast"
import "go/types"
import "text/template"

func init() {
	vd := &VerbDefinition{Verb: "index"}
	vd.Description = "returns the element of the specified slice valued field at the specified (zero based) index."
	vd.Assimilate = func(ctx *context, vd *VerbDefinition, spec *slotSpec, id *InterfaceDefinition, m *ast.Field) error {
		ftype, ok := m.Type.(*ast.FuncType)
		if !ok {
			return nil
		}
		if err := checkSignature(ftype, 1, 1); err != nil {
			return err
		}
		params := util.FieldListSlice(ftype.Params)
		if types.ExprString(params[0].Type) != "int" {
			return fmt.Errorf("The parameter of an 'index' method must be type int")
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
		func (x *{{.StructName}}) {{.MethodName}} (index int) {{ExprString .Type.Elt}} {
			return x.{{.SlotName}}[index]
		}
	`))
	VerbDefinitions[vd.Verb] = vd
}

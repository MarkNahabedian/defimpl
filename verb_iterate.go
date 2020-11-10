package main

import "fmt"
import "go/ast"
import "go/types"
import "text/template"
import "defimpl/util"

func init() {
	vd := &VerbDefinition{
		Verb:         "iterate",
		ParamCount:   1,
	}
	vd.Description = "Applies the specified function to each element of the slice-valued slot until the function returns false."
	vd.Assimilate = func(ctx *context, vd *VerbDefinition, spec *slotSpec, id *InterfaceDefinition, m *ast.Field) error {
		ftype, ok := m.Type.(*ast.FuncType)
		if !ok {
			return nil
		}
		if err := checkSignature(ftype, 1, 0); err != nil {
			return err
		}
		params := util.FieldListSlice(ftype.Params)

		// params[0] should be of type func(spec.Type) bool
		funarg, ok := params[0].Type.(*ast.FuncType)
		if !ok {
			return fmt.Errorf("Wrong type for iterate parameter")
		}
		if err := checkSignature(funarg, 1, 1); err != nil {
			return err
		}
		/*
			if ctx.info.TypeOf(funarg.Results.List[0].Type).String() != "bool" {
				return fmt.Errorf("Wrong return type for iterate parameter")
			}
		*/
		spec.CheckType(SliceOfType(funarg.Params.List[0].Type))
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
	vd.TopLevelTemplate = template.Must(template.New(vd.Verb).Funcs(map[string]interface{}{
		"ExprString": types.ExprString,
	}).Parse(`
		{{.DocComment}}
		func (x *{{.StructName}}) {{.MethodName}} (f func(item {{ExprString .Type.Elt}}) bool) {
			for _, v := range x.{{.SlotName}} {
				if !f(v) {
					break
				}
			}
		}
	`))
	VerbDefinitions[vd.Verb] = vd
}

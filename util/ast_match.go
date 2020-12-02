package util

import "fmt"
import "reflect"
import "strings"
import "go/ast"
import "go/token"


// AstMatch matches a candidate ast.Node against a pattern ast.Node.
// If the name of an ast.Ident in pattern starts with '_', then that
// name and the corresponding value from candidate are added to
// scratchpad.
func AstMatch(pattern interface{}, candidate interface{}, scratchpad map[string]interface{}) (bool, error) {
	if reflect.TypeOf(pattern) != reflect.TypeOf(candidate) {
		return false, fmt.Errorf("type mismatch: %T, %T", pattern, candidate)
	}
	switch p := pattern.(type) {
	case token.Pos:
		// We don't care whether positions match.
		return true, nil
	case *ast.Ident:
		c := candidate.(*ast.Ident)
		if strings.HasPrefix(p.Name, "_") {
			val, ok := scratchpad[p.Name]
			if ok {
				if name, ok := val.(string); ok {
					if name != c.Name {
						return false, fmt.Errorf("identifier mismatch: %s, %s", name, c.Name)
					}
				} else {
					return false, fmt.Errorf("scratchpad[%s] is %T, not string", p.Name, val)
				}
			} else {
				scratchpad[p.Name] = c
			}
		}
	default:
		pv := reflect.ValueOf(pattern)
		cv := reflect.ValueOf(candidate)
		if pv.Kind() == reflect.Array || pv.Kind() == reflect.Slice {
			if pvl, cvl :=  pv.Len(), cv.Len(); pvl != cvl {
				return false, fmt.Errorf("Lengths don't match: %d, %d", pvl, cvl)
			} else {
				for i := 0; i < pvl; i++ {
					match, err := AstMatch(
						pv.Index(i).Interface(),
						cv.Index(i).Interface(),
						scratchpad)
					if !match {
						return false, err
					}
				}
			}
			return true, nil
		}
		
		if pv.Kind() != reflect.Ptr {
			panic(fmt.Errorf("%T %T not a pointer", pattern, candidate))
		}
		pvi := reflect.Indirect(pv)
		cvi := reflect.Indirect(cv)
		if pvi.Kind() == reflect.Invalid && cvi.Kind() == reflect.Invalid {
			// Null pointers match.
			return true, nil
		}
		if pvi.Kind() != reflect.Struct || cvi.Kind() != reflect.Struct {
			panic(fmt.Errorf("%T, %T not struct pointer", pattern, candidate))
		}
		for i := 0; i < reflect.Indirect(pv).NumField(); i++ {
			pf := pvi.Field(i).Interface()
			cf := cvi.Field(i).Interface()
			// fmt.Printf("Considering %v %v\n", pf, cf)
			if match, err := AstMatch(pf, cf, scratchpad); !match {
				return false, err
			}
		}
	}
	return true, nil
}


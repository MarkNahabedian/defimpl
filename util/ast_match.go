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
	switch p := pattern.(type) {
	case token.Pos:
		// We don't care whether positions match.
		return true, nil
	case *ast.Ident:
		if strings.HasPrefix(p.Name, "_") {
			_, ok := scratchpad[p.Name]
			if ok {
				return false, fmt.Errorf("%s already set", p.Name)
			} else {
				scratchpad[p.Name] = candidate
			}
		}
	case *ast.Field:
		if reflect.TypeOf(pattern) != reflect.TypeOf(candidate) {
			return false, fmt.Errorf("type mismatch: %T, %T", pattern, candidate)
		}
		return AstMatch(p.Type, candidate.(*ast.Field).Type, scratchpad)
	default:
		// Kludge because I don't have the patience to
		// implement a clause for every type we might
		// encounter in an ast.
		if reflect.TypeOf(pattern) != reflect.TypeOf(candidate) {
			return false, fmt.Errorf("type mismatch: %T, %T", pattern, candidate)
		}
		pv := reflect.ValueOf(pattern)
		cv := reflect.ValueOf(candidate)
		if pv.Kind() == reflect.Array || pv.Kind() == reflect.Slice {
			if pvl, cvl :=  pv.Len(), cv.Len(); pvl != cvl {
				return false, fmt.Errorf("Lengths don't match: %d, %d; %v, %v",
					pvl, cvl,
					pv.Interface(), cv.Interface())
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
			if match, err := AstMatch(pf, cf, scratchpad); !match {
				return false, err
			}
		}
	}
	return true, nil
}


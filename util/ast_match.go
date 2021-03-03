package util

import "fmt"
import "reflect"
import "strings"
import "go/ast"
import "go/token"


var AST_MATCH_DEBUG_DUMP = false
var AST_MATCH_DEBUG_DUMP_ALL = false


// AstMatch matches a candidate ast.Node against a pattern ast.Node.
//
// If the name of an ast.Ident in pattern starts with '_', then that
// name and the corresponding value from candidate are added to
// scratchpad.
//
// If IGNORE appears as an identifier in pattern then it matches the
// corresponding element of candidate.
func AstMatch(pattern interface{}, candidate interface{}, scratchpad map[string]interface{}) (bool, error) {
	const err_prefix = "defimpl/util:AstMatch: "
	type stackLevel struct {
		pattern interface{}
		candidate interface{}
		previous *stackLevel
	}
	// Dump the pattern and candidate for debugging:
	dump := func(stack *stackLevel) {
		if AST_MATCH_DEBUG_DUMP {
			fmt.Printf("AstMatch stack:\n")
			for s := stack; s != nil; s = s.previous {
				if AST_MATCH_DEBUG_DUMP_ALL || s.previous == nil {
					fmt.Printf("\n")
					ast.Print(nil, s.pattern)
					ast.Print(nil, s.candidate)
				}
			}
		}
	}
	var astm func (pattern, candidate interface{}, stack *stackLevel) (bool, error)
	astm = func (pattern, candidate interface{}, stack *stackLevel) (bool, error) {
		stack = &stackLevel {
			pattern: pattern,
			candidate: candidate,
			previous: stack,
		}
		if pattern == nil {
			if candidate == nil {
				return true, nil
			} else {
				return false, fmt.Errorf("%spattern is nil but candidate %#v isn't",
					err_prefix, candidate)
			}
		}
		switch p := pattern.(type) {
		case token.Pos:
			// We don't care whether positions match.
			return true, nil
		case *ast.Ident:
			if p.Name == "IGNORE" {
				return true, nil
			}
			if strings.HasPrefix(p.Name, "_") {
				_, ok := scratchpad[p.Name]
				if ok {
					return false, fmt.Errorf("%s%s already set",
						err_prefix, p.Name)
				} else {
					scratchpad[p.Name] = candidate
					return true, nil
				}
			}
			if c, ok := candidate.(*ast.Ident); ok {
				if p.Name == c.Name {
					return true, nil
				}
				return false, fmt.Errorf("%sIdent Names don't match: %s, %s",
					err_prefix, p.Name, c.Name)
			}
			return false, fmt.Errorf("%sExpected candidate to be ast.Ident, not %T",
				err_prefix, candidate)
		case *ast.Field:
			if reflect.TypeOf(pattern) != reflect.TypeOf(candidate) {
				return false, fmt.Errorf("%stype mismatch: %T, %T",
					err_prefix, pattern, candidate)
			}
			return astm(p.Type, candidate.(*ast.Field).Type, stack)
		case *ast.FieldList:
			if ok, name := wholeFieldList(p); ok {
				scratchpad[name] = candidate
				return true, nil
			}
			if c, ok := candidate.(*ast.FieldList); ok {
				if c == nil {
					if p == nil || p.List == nil || len(p.List) == 0 {
						return true, nil
					} else {
						return false, fmt.Errorf("%scandidate is nil against non-empty FieldList pattern",
							err_prefix)
					}
				}
				if len(p.List) != len(c.List) {
					return false, fmt.Errorf("%sFieldLists differ in length",
						err_prefix)
				}
				for i := 0; i < len(p.List); i++ {
					if ok, err := astm (p.List[i], c.List[i], stack); !ok {
						return false, err
					}
				}
				return true, nil
			} else {
				return false, fmt.Errorf("%scandidate isn't FieldList", err_prefix)
			}
		}
		// All of the switch clauses above are expected to
		// have returned with an answer unless they want to
		// drop through to this default case. This is not a
		// default clause of the switch in case other clauses
		// want to delegate to it.
		//
		// Kludge because I don't have the patience to
		// implement a clause for every type we might
		// encounter in an ast.
		if reflect.TypeOf(pattern) != reflect.TypeOf(candidate) {
			return false, fmt.Errorf("%stype mismatch: %T, %T",
				err_prefix, pattern, candidate)
		}
		pv := reflect.ValueOf(pattern)
		cv := reflect.ValueOf(candidate)
		if ((pv.Kind() == reflect.Array || pv.Kind() == reflect.Slice) &&
			(cv.Kind() == reflect.Array || cv.Kind() == reflect.Slice)) {
			if pvl, cvl :=  pv.Len(), cv.Len(); pvl != cvl {
				return false, fmt.Errorf("%sLengths don't match: %d, %d; %v, %v",
					err_prefix, pvl, cvl,
					pv.Interface(), cv.Interface())
			} else {
				for i := 0; i < pvl; i++ {
					match, err := astm(
						pv.Index(i).Interface(),
						cv.Index(i).Interface(),
						stack)
					if !match {
						return false, err
					}
				}
			}
			return true, nil
		}
		if pv.Kind() != reflect.Ptr {
			dump(stack)
			panic(fmt.Errorf("%s%T %T not a pointer",
				err_prefix, pattern, candidate))
		}
		pvi := pv.Elem()
		cvi := cv.Elem()
		if pvi.Kind() == reflect.Invalid && cvi.Kind() == reflect.Invalid {
			// Null pointers match.
			return true, nil
		}
		if pvi.Kind() != reflect.Struct || cvi.Kind() != reflect.Struct {
			dump(stack)
			panic(fmt.Errorf("%s%T (kind:%s), %T (kind:%s) not struct pointer",
				err_prefix, pattern, pvi.Kind(),
				candidate, cvi.Kind()))
		}
		if pvi.Type() != cvi.Type() {
			return false, fmt.Errorf("%sdifferent struct types %v, %v",
				err_prefix, pvi.Type(), cvi.Type())
		}
		for i := 0; i < reflect.Indirect(pv).NumField(); i++ {
			pf := pvi.Field(i).Interface()
			cf := cvi.Field(i).Interface()
			if match, err := astm(pf, cf, stack); !match {
				return false, err
			}
		}
		return true, nil
	}
	return astm(pattern, candidate, (*stackLevel)(nil))
}

func wholeFieldList(f *ast.FieldList) (bool, string) {
	s := FieldListSlice(f)
	if len(s) != 1 {
		return false, ""
	}
	i, ok := s[0].Type.(*ast.Ident)
	if !ok {
		return false, ""
	}
	if strings.HasPrefix(i.Name, "__") {
		return true, i.Name
	}
	return false, ""	
}

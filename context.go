package main

import "fmt"
import "regexp"
import "go/ast"
import "go/importer"
import "go/parser"
import "go/token"
import "go/types"
import "os"
import "path/filepath"


// context is the top level oblect representing the task of running
// defimpl for a single go package source directory.
type context struct {
	dir   string
	fset  *token.FileSet
	info  *types.Info
	astFiles []*ast.File
	files []*File
	typeErrors []error
}

// NewContext returns a context for orchestrating defimpl's operations.
// dir should be an absolute path to a go package source directory.
// The go source files in dir will be parsed and File objects added to
// the files field of the new context.
func NewContext(dir string) (*context, error) {
	if !filepath.IsAbs(dir) {
		return nil, fmt.Errorf("%s is not an absolute path", dir)
	}
	ctx := &context{dir: dir}
	ctx.fset = token.NewFileSet()
	ctx.info = &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}
	pkgs, err := parser.ParseDir(ctx.fset, dir,
		func(fi os.FileInfo) bool {
			return !IsOutputFilePath(fi.Name())
		}, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	for _, pkg := range pkgs {
		for _, astFile := range pkg.Files {
			ctx.astFiles = append(ctx.astFiles, astFile)
		}
	}
	ctx.Check()
	for _, astFile := range ctx.astFiles {
		// NewFile is where interface definitions are
		// processed and VerbPhrases created.
		ctx.files = append(ctx.files, NewFile(ctx, astFile))
	}
	ctx.ReportTypeErrors()
	return ctx, nil
}

// Check runs the go type checker on all of the files in ctx.
func (ctx *context) Check() {
	conf := types.Config{
		Importer: importer.For("source", nil), // importer.Default(),
		Error: func(err error) {
			ctx.typeErrors = append(ctx.typeErrors, err)
		},
	}
	_, _ = conf.Check(ctx.astFiles[0].Name.Name, ctx.fset, ctx.astFiles, ctx.info)
}

func (ctx *context) ReportTypeErrors() {
	ignore := false
	for _, err := range ctx.typeErrors {
		missing := typeErrorUndeclaredName(err)
		if missing != "" {
			for _, f := range ctx.files {
				for _, idef := range f.Interfaces {
					if idef.StructName() == missing {
						ignore = true
						break 
					}
				}
				if ignore {
					break
				}
			}
		}
		if !ignore {
			fmt.Fprintf(os.Stderr, "defimpl error while type checking: %s\n", err)
		}
	}
}

var typeErrorUndeclaredNameRegexp = regexp.MustCompile(
	`^undeclared name: (?P<type>[a-zA-Z_0-9]+)$`)

func typeErrorUndeclaredName(err error) string {
	e, ok := err.(types.Error)
	if !ok {
		return ""
	}
	m := typeErrorUndeclaredNameRegexp.FindStringSubmatch(e.Msg)
	if len(m) > 1 {
		return m[1]
	}
	return ""
}

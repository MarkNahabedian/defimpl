package main

import "fmt"
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
	dir string
	fset *token.FileSet
	info *types.Info
	files []*File
}

// NewContext returns a context for orchestrating defimpl's operations.
// dir should be an absolute path to a go package source directory.
// The go source files in dir will be parsed and File objects added to
// the files field of the new context.
// context.Check should be run separately.
func NewContext(dir string) (*context, error) {
	if !filepath.IsAbs(dir) {
		return nil, fmt.Errorf("%s is not an absolute path", dir)
	}
	ctx := &context{ dir: dir }
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
			ctx.files = append(ctx.files, NewFile(ctx, astFile))
		}
	}
	return ctx, nil
}

// Check runs the go type checker on all of the files in ctx.
func (ctx *context) Check() error {
	astFiles := []*ast.File {}
	for _, f := range ctx.files {
		astFiles = append(astFiles, f.AstFile)
	}
	conf := types.Config{
		Importer: importer.For("source", nil), // importer.Default(),
	}
	_, err := conf.Check(astFiles[0].Name.Name, ctx.fset, astFiles, ctx.info)
	if err != nil {
		return err
	}
	return nil
}

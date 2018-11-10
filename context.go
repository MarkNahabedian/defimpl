package main

import "fmt"
import "go/parser"
import "go/token"
import "os"
import "path/filepath"

// context is the top level oblect representing the task of running
// defimpl for a single go package source directory.
type context struct {
	dir   string
	fset  *token.FileSet
//	info  *types.Info
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
	ctx := &context{dir: dir}
	ctx.fset = token.NewFileSet()
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

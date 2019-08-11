package main

import "bytes"
import "fmt"
import "go/ast"
import "go/format"
import "go/parser"
import "go/types"
import "os"
import "path/filepath"
import "strings"
import "text/template"

// File represents a single file to be processed.
type File struct {
	AstFile       *ast.File
	InputFilePath string
	Package       string
	Interfaces    []*InterfaceDefinition
}

func (f *File) Defimpl() string {
	return os.Args[0]
}

func (f *File) OutputFilePath() string {
	input := f.InputFilePath
	return filepath.Join(filepath.Dir(input), "impl_"+filepath.Base(input))
}

func IsOutputFilePath(f string) bool {
	return strings.HasPrefix(filepath.Base(f), "impl_")
}

// AddImports addresses the need to include imports for captured
// refernces like type definitions.
// It will need to be used for each types.Type that is included in the output file.
func (f *File) AddImports(ctx *context, astFile *ast.File) {
	for _, i := range f.Interfaces {
		i.AddImports(ctx, f.AstFile, astFile)
	}
}

// NewFile returns a File object for the given ast.File.
// Interfaces will be filled in.
func NewFile(ctx *context, astFile *ast.File) *File {
	f := &File{
		AstFile:       astFile,
		InputFilePath: ctx.fset.Position(astFile.Package).Filename,
		Package:       astFile.Name.Name,
		Interfaces:    []*InterfaceDefinition{},
	}
	for _, decl := range astFile.Decls {
		if id := NewInterface(ctx, f.Package, decl); id != nil {
			f.Interfaces = append(f.Interfaces, id)
		}
	}
	return f
}

func (f *File) Write(ctx *context) error {
	if !f.AnyStructs() {
		return nil
	}
	output := f.OutputFilePath()
	out, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("Can't create %s: %s", output, err)
	}
	format.Node(out, ctx.fset, f.GenerateCode(ctx))
	out.Close()
	fmt.Printf("Wrote %s\n", output)
	return nil
}

func (f *File) AnyStructs() bool {
	for _, i := range f.Interfaces {
		if i.DefinesStruct() {
			return true
		}
	}
	return false
}

func (f *File) GenerateCode(ctx *context) *ast.File {
	writer := bytes.NewBufferString("")
	err := OutputFileTemplate.Execute(writer, f)
	if err != nil {
		panic(err)
	}
	parsed, err := parser.ParseFile(ctx.fset, "", writer.String(), parser.ParseComments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", writer.String())
		panic(err)
	}
	f.AddImports(ctx, parsed)
	return parsed
}

// OutputFileTemplate is the template for generating the output file
// containing the programatically generated struct and method definitions
// that implement the interfaces in the input file.
// The parameter is a File object.
var OutputFileTemplate *template.Template = template.Must(template.New("OutputFileTemplate").Funcs(map[string]interface{}{
	//		"NormalizedType": util.NormalizedType,
	"ExprString": types.ExprString,
}).Parse(`
// This file was automatically generated by {{.Defimpl}} from {{.InputFilePath}}.
package {{.Package}}
	{{with $file := .}}
		{{range .Interfaces}}
			{{if .DefinesStruct}}
				{{if not .IsAbstract}}
					type {{.StructName}} struct {
						{{range .SlotSpecs}}	
						{{.Name}} {{ExprString .Type}}
					 	{{end}}
						{{- /* Fields required to support abstract inherited interfaces: */ -}}
						{{with $thisInterface := .}}
							{{range $inherited := .AllInherited}}
								{{if $inherited.IsAbstract}}
									// Fields to support the {{$inherited.InterfaceName}} interface:
									{{range $inherited.SlotSpecs}}
										{{.Name}} {{ExprString .Type}}
								 	{{end}}
								{{else}}
									// Interface {{$inherited.InterfaceName}} has a concrete implementation
									{{$inherited.StructName}}
								{{end}}
							{{end}}
						{{end}}
					}
					{{with $interface := .}}
						{{range .SlotSpecs}}
							{{range .Verbs}}
								{{.RunTemplate}}
							{{end}}
						{{end}}
						{{range .InheritedVerbs}}
							{{.RunTemplate}}
						{{end}}
					{{end}}
				{{end}}
			{{end}}
		{{end}}
	{{end}}
`)) // end template

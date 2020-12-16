package main

import "flag"
import "fmt"
import "os"

var debug_dump = false
func init() {
	flag.BoolVar(&debug_dump, "debug_dump", false,
		"Dump internal data when debugging defimpl itself.")
}

func (ctx *context) debug_dump() {
	for _, file := range ctx.files {
		file.debug_dump()
	}
}

func (file *File) debug_dump() {
	fmt.Fprintf(os.Stderr, "File %s\n", file.InputFilePath	)
	for _, idef := range file.Interfaces {
		idef.debug_dump()
	}
}

func (idef *InterfaceDefinition) debug_dump() {
	abstract := ""
	if idef.IsAbstract {
		abstract = "abstract "
	}
	fmt.Fprintf(os.Stderr, "  interface %s.%s%s:\n",
		idef.Package(), idef.InterfaceName, abstract)
	for _, vp := range idef.VerbPhrases {
		debug_dump_vp(vp)
	}
}


func debug_dump_vp(vp VerbPhrase) {
	// Unsafe.  StructBody has side effects.

	sb := func() string {
		s, err := vp.Verb().StructBody(vp)
		if err != nil {
			return err.Error()
		}
		return s
	}
	sb()
	switch vpt := vp.(type) {
	case SlotVerbPhrase:
		fmt.Fprintf(os.Stderr, "    %s %s %s %v\n",
			vpt.Verb().Tag(), vpt.MethodName(),
			vpt.SlotName(), vpt.SlotType(),
			// sb()
		)
	default:
		fmt.Fprintf(os.Stderr, "    %s %s\n",
			vpt.Verb().Tag(), vpt.MethodName(),
			// sb()
		)
	}
}



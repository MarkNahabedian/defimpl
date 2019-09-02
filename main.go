package main

import "flag"
import "fmt"
import "os"
import "path/filepath"

var show_verbs bool = false
func init() {
	flag.BoolVar(&show_verbs, "show_verbs", false, "Just list supported defimpl verbs and exit.")
}

func main() {
	flag.Parse()
	if show_verbs {
		for _, v := range VerbDefinitions {
			fmt.Fprintf(os.Stderr, "%s\t  %s\n",
				v.Verb, v.Description)
		}
		return
	}
	afp, err := filepath.Abs(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't determine working directory: %s\n", err)
		return
	}
	ctx, err := NewContext(filepath.Clean(afp))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}
	ctx.DoInheritance()
	for _, f := range ctx.files {
		fmt.Printf("file %s\n", f.InputFilePath)
		if err := f.Write(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
	}
}

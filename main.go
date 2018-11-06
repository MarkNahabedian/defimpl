package main

import "fmt"
import "os"
import "path/filepath"

func main() {
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
	/*
		if err := ctx.Check(); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return
		}
	*/
	for _, f := range ctx.files {
		fmt.Printf("file %s\n", f.InputFilePath)
		if err := f.Write(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
	}
}

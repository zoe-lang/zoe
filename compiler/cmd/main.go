package main

import (
	"fmt"
	"log"
	"os"

	zoe "github.com/ceymard/zoe/compiler"
	"github.com/fatih/color"
)

// https://github.com/sourcegraph/go-lsp
var blue = color.New(color.FgHiBlue).SprintFunc()
var yel = color.New(color.FgHiYellow).SprintFunc()

func main() {

	for _, fname := range os.Args[1:] {

		_, _ = fmt.Print("\nHandling ", yel(fname), "\n")
		file, err := zoe.NewFile(fname)
		if err != nil {
			log.Printf("-- %v", err)
			// if ctx != nil && ctx.Start != nil {
			// 	_, _ = fmt.Print(ctx.Start.ToSlice())
			// }
			continue
		}

		file.Parse()
		// log.Print(file.Nodes)
		file.PrintNode(os.Stderr, file.RootNode)
		_, _ = os.Stderr.WriteString("\n")
		file.TestFileAst()
		// log.Printf("%v", file.Nodes)
		// _, _ = os.Stdout.WriteString(res.DumpString() + "\n\n")
		// file.TestFileAst()

		for _, err := range file.Errors {
			err.Print(os.Stderr)
		}
	}
}

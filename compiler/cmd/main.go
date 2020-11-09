package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	zoe "github.com/ceymard/zoe/compiler"
	"github.com/fatih/color"
)

// https://github.com/sourcegraph/go-lsp
var blue = color.New(color.FgHiBlue).SprintFunc()

func main() {

	for _, fname := range os.Args[1:] {

		fmt.Print(`Handling '`, blue(fname), "'\n\n")
		ctx, err := zoe.NewZoeContext(fname)
		if err != nil {
			log.Print(err)
			continue
		}

		res := ctx.ParseFile()
		var buf bytes.Buffer
		res.Dump(&buf)
		_, _ = os.Stdout.WriteString(buf.String() + "\n\n")
		ctx.TestFileAst()

		for _, err := range ctx.Errors {
			err.Print(os.Stderr)
		}
	}
}

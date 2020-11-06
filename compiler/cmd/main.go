package main

import (
	"fmt"
	"log"
	"os"

	zoe "github.com/ceymard/zoe/compiler"
	"github.com/fatih/color"
)

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
		_, _ = os.Stdout.WriteString(res.Debug() + "\n\n")
		ctx.TestFileAst()

		for _, err := range ctx.Errors {
			err.Print(os.Stderr)
		}
	}
}

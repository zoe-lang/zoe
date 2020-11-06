package main

import (
	"log"
	"os"

	zoe "github.com/ceymard/zoe/compiler"
	"github.com/fatih/color"
)

var blue = color.New(color.FgHiBlue).SprintFunc()

func main() {

	for _, fname := range os.Args[1:] {

		log.Print(`Handling '`, blue(fname), `'`)
		ctx, err := zoe.NewZoeContext(fname)
		if err != nil {
			log.Print(err)
			continue
		}

		res := ctx.ParseFile()
		_, _ = os.Stdout.WriteString(res.Debug() + "\n")
		ctx.TestFileAst()

		for _, err := range ctx.Errors {
			err.Print(os.Stderr)
		}
	}
}

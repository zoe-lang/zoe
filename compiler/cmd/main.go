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

func handleFile(fname string) *zoe.File {
	_, _ = fmt.Print("Handling ", yel(fname), "\n")
	file, err := zoe.NewFile(fname)
	if err != nil {
		log.Printf("-- %v", err)
		// if ctx != nil && ctx.Start != nil {
		// 	_, _ = fmt.Print(ctx.Start.ToSlice())
		// }
		return nil
	}

	file.Parse()

	// os.Stderr.Write([]byte("Pouet pouet !\n"))
	if err != nil {
		_, _ = os.Stderr.Write([]byte(err.Error() + "\n"))
	} else {
		// if res, err := json.Marshal(file.RootNode); err != nil {
		// 	_, _ = os.Stderr.Write([]byte("json: "))
		// 	_, _ = os.Stderr.Write(res)
		// }
	}

	_, _ = os.Stderr.WriteString("\n")
	file.TestFileAst()
	for _, err := range file.Errors {
		err.Print(os.Stderr)
	}
	return file
}

func main() {

	for _, fname := range os.Args[1:] {
		handleFile(fname)
	}
}

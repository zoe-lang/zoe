package zoe

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func handleFile(fname string) *File {
	_, _ = fmt.Print("Handling ", yel(fname), "\n")
	file, err := NewFile(fname)
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

	// _, _ = os.Stderr.WriteString("\n")
	// file.TestFileAst()
	for _, err := range file.Errors {
		err.Print(os.Stderr)
	}
	return file
}

func TestFiles(t *testing.T) {
	var total = 0
	var handle = func(path string, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, ".zo") {
			var file = handleFile(path)
			if file != nil {
				total += len(file.Errors)
			}
		}
		return nil
	}
	var handleDir = func(path string) {
		_ = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			return handle(path, err)
		})
	}
	handleDir("../tests")
	handleDir("../../std")

	log.Print("  --> total errors : ", total)
}

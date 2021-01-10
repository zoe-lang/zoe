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
	_, _ = fmt.Print("\nHandling ", yel(fname), "\n")
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

	_, _ = os.Stderr.WriteString("\n")
	file.TestFileAst()
	for _, err := range file.Errors {
		err.Print(os.Stderr)
	}
	return file
}

func TestFiles(t *testing.T) {
	var total = 0
	_ = filepath.Walk("../tests", func(path string, info os.FileInfo, err error) error {
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
	})
	log.Print("  --> total errors : ", total)
}

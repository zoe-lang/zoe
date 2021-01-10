package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

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

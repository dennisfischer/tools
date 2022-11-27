package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/dennisfischer/tools/indexer"
)

const (
	rootDir    = "/mnt/e/"
	outputPath = "/mnt/d/files.binproto"
)

func main() {
	output, err := os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("could not open output file: %v", err)
	}
	defer output.Close()

	idx := indexer.NewIndexer()
	idx.Load(output)

	filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err.Error())
		}
		if info.IsDir() {
			return nil
		}
		idx.AddFile(output, path)
		return nil
	})
}

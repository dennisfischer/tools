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
	outputPath = "/mnt/d/files_backup_disk.binproto"
)

func main() {
	input, err := os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("could not open output file: %v", err)
	}
	defer input.Close()

	oldIdx := indexer.NewIndexer()
	oldIdx.Load(input)
	log.Printf("Restored index with %d files.", len(oldIdx.Index.FilesByPath))

	output, err := os.CreateTemp(filepath.Dir(outputPath), filepath.Base(outputPath))
	if err != nil {
		log.Fatalf("could not create temporary output file: %v", err)
	}
	defer os.Remove(output.Name())

	newIdx := indexer.NewIndexer()
	filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err.Error())
		}
		if info.IsDir() {
			return nil
		}
		if existing, ok := oldIdx.Index.FilesByPath[path]; ok {
			newIdx.AddFileNoIndex(output, existing)
		} else {
			newIdx.AddFile(output, path)
		}
		return nil
	})

	err = os.Rename(output.Name(), outputPath)
	if err != nil {
		log.Fatalf("could not rename temporary file: %v", err)
	}
}

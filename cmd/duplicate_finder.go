package main

import (
	"log"
	"os"

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

	log.Printf("Restored index with %d files.", len(idx.Index.FilesByPath))
	for _, v := range idx.Index.FilesByHash {
		if len(v) > 1 {
			log.Println(v)
		}
	}
}

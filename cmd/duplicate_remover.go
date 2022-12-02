package main

import (
	"log"
	"os"
	"sort"
	"strings"

	"github.com/dennisfischer/tools/indexer"
)

const (
	backupDir  = "/mnt/e/backup_disk/"
	outputPath = "/mnt/d/files_backup_disk.binproto"
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
			sort.SliceStable(v, func(i, j int) bool {
				if strings.HasPrefix(*v[i].Path, backupDir) {
					if strings.HasPrefix(*v[j].Path, backupDir) {
						return *v[i].Path < *v[j].Path
					}
					return false
				}
				if strings.HasPrefix(*v[j].Path, backupDir) {
					return true
				}

				return *v[i].Path < *v[j].Path
			})

			keep := v[0]
			if _, err := os.Stat(*keep.Path); err != nil {
				log.Printf("file to keep is missing: %v\n", err)
				continue
			}
			log.Printf("keeping file: %s\n", *keep.Path)

			others := v[1:]
			for _, f := range others {
				log.Printf("removing file: %s\n", *f.Path)
				err = os.Remove(*f.Path)
				if err != nil {
					log.Fatalf("cannot remove file: %v\n", err)
				}
			}
		}
	}
}

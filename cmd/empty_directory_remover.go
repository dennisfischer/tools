package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

const (
	rootDir = "/mnt/e/"
)

func main() {
	filepath.Walk(rootDir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err.Error())
		}
		if !info.IsDir() {
			return nil
		}

		if empty, err := IsEmpty(p); err != nil {
			log.Printf("cannot remove directory: %v\n", err)
		} else if empty {
			log.Printf("removing empty directory: %v\n", p)
			os.Remove(p)
		}
		return nil
	})
}

func IsEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

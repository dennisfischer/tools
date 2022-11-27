package indexer

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"google.golang.org/protobuf/proto"

	fp "github.com/dennisfischer/tools/proto"
)

type Index struct {
	Files map[string]fp.File
}

type Indexer struct {
	index Index
}

func NewIndexer() Indexer {
	return Indexer{
		index: Index{
			Files: make(map[string]fp.File, 0),
		},
	}
}

func (i *Indexer) Load(r io.Reader) {
	for {
		b, err := read(r)
		if err != nil {
			if err == io.EOF { // Reached end of input
				break
			}
			log.Fatalf("could not read protobuf from file: %v", err)
		}
		f := &fp.File{}
		err = proto.Unmarshal(b, f)
		if err != nil {
			log.Fatalf("could not unmarshal protobuf: %v", err)
		}

		i.index.Files[*f.Path] = *f
	}
}

func (i *Indexer) WriteAll(w io.Writer) {
}

func (i *Indexer) AddFile(w io.Writer, path string) error {
	if _, ok := i.index.Files[path]; ok {
		log.Printf("Skipping file: %s\n", path)
		return nil
	}

	log.Printf("File Name: %s\n", path)
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := hashFile(f)
	n := filepath.Base(path)
	t := filepath.Ext(n)

	file := fp.File{
		Name: &n,
		Path: &path,
		Hash: &h,
		Type: &t,
	}
	err = write(w, &file)
	if err != nil {
		log.Fatalln(err)
	}

	i.index.Files[path] = file
	return nil
}

func hashFile(r io.Reader) string {
	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(h.Sum(nil))
}

func read(r io.Reader) ([]byte, error) {
	buf := make([]byte, 4)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, err
	}

	size := binary.LittleEndian.Uint32(buf)

	msg := make([]byte, size)
	if _, err := io.ReadFull(r, msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func write(w io.Writer, m proto.Message) error {
	msg, err := proto.Marshal(m)
	if err != nil {
		return fmt.Errorf("could not marshal proto: %v", err)
	}
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(len(msg)))

	if _, err := w.Write(buf); err != nil {
		return fmt.Errorf("could not write proto size to file: %v", err)
	}

	if _, err := w.Write(msg); err != nil {
		return fmt.Errorf("could not write proto to file: %v", err)
	}
	return nil
}

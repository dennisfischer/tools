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
	FilesByPath map[string]*fp.File
	FilesByHash map[string][]*fp.File
}

type Indexer struct {
	Index *Index
}

func NewIndexer() Indexer {
	return Indexer{
		Index: &Index{
			FilesByPath: make(map[string]*fp.File, 0),
			FilesByHash: make(map[string][]*fp.File, 0),
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

		i.Index.addFileToIndex(f)
	}
}

func (i *Indexer) WriteAll(w io.Writer) {
}

func (i *Indexer) AddFile(w io.Writer, path string) error {
	if _, ok := i.Index.FilesByPath[path]; ok {
		return nil
	}

	log.Printf("File Name: %s\n", path)
	fh, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()

	h := hashFile(fh)
	n := filepath.Base(path)
	t := filepath.Ext(n)

	f := &fp.File{
		Name: &n,
		Path: &path,
		Hash: &h,
		Type: &t,
	}
	i.AddFileNoIndex(w, f)
	return nil
}

func (i *Indexer) AddFileNoIndex(w io.Writer, f *fp.File) error {
	err := write(w, f)
	if err != nil {
		log.Fatalln(err)
	}

	i.Index.addFileToIndex(f)
	return nil
}

func (i *Index) addFileToIndex(f *fp.File) {
	i.FilesByPath[*f.Path] = f
	if _, ok := i.FilesByHash[*f.Hash]; !ok {
		i.FilesByHash[*f.Hash] = []*fp.File{f}
	} else {
		i.FilesByHash[*f.Hash] = append(i.FilesByHash[*f.Hash], f)
	}
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

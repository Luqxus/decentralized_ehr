package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const defaultRootFolder = "spxce"

// PathTransformFunc transform the file path
// created from the file content hash sum
// return PathKey
type PathTransformFunc func(string) PathKey

type StoreOpts struct {
	// Root is the name of the root folder
	// contains all the folders/files of the system
	Root string

	// PathTransformFunc transform the file path
	// created from the file content hash sum
	// return PathKey
	PathTransformFunc PathTransformFunc
}

// DefaultPathTransformFunc is used if no custom transform is provided
var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		Pathname: key,
		Filename: key,
		Root:     key,
	}
}

type Store struct {
	StoreOpts
}

// A container of the paths and filename of a file
type PathKey struct {
	// the pathname to the file
	Pathname string
	// the name of the file
	Filename string
	// the root folder of the path
	Root string
}

func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))

	hashStr := hex.EncodeToString(hash[:])

	blockSize := 5

	sliceLen := len(hashStr) / blockSize

	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i*blockSize)+blockSize
		paths[i] = hashStr[from:to]
	}

	return PathKey{
		Pathname: strings.Join(paths, "/"),
		Filename: hashStr,
		Root:     hashStr[:blockSize],
	}
}

func NewStore(opts StoreOpts) *Store {
	if opts.PathTransformFunc == nil {
		opts.PathTransformFunc = DefaultPathTransformFunc
	}

	if len(opts.Root) == 0 {
		opts.Root = defaultRootFolder
	}

	return &Store{
		StoreOpts: opts,
	}
}

func (p *PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s", p.Pathname, p.Filename)
}

func (s *Store) Has(key string) bool {
	pathKey := s.PathTransformFunc(key)

	_, err := os.Stat(s.Root + "/" + pathKey.FullPath())

	// @luqxus check for error properly
	return errors.Is(err, nil)
}

// clears all system files
func (s *Store) Clear() error {
	return os.RemoveAll(s.Root)
}

func (s *Store) Delete(key string) error {
	pathKey := s.PathTransformFunc(key)

	defer func() {
		fmt.Printf("deleted [%s] from disk\n", pathKey.Filename)
	}()

	return os.RemoveAll(s.Root + "/" + pathKey.Root)
}

func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	buf := new(bytes.Buffer)

	io.Copy(buf, f)

	return buf, nil
}

func (s *Store) Write(key string, r io.Reader) (int64, error) {
	return s.writeStream(key, r)
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(key)

	return os.Open(s.Root + "/" + pathKey.FullPath())

}

func (s *Store) writeStream(key string, r io.Reader) (int64, error) {
	pathKey := s.PathTransformFunc(key)

	if err := os.MkdirAll(s.Root+"/"+pathKey.Pathname, os.ModePerm); err != nil {
		return 0, err
	}

	f, err := os.Create(s.Root + "/" + pathKey.FullPath())
	if err != nil {
		return 0, err
	}

	n, err := io.Copy(f, r)
	if err != nil {
		return 0, err
	}

	return n, nil
}

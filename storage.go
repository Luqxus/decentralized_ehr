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

		// path to file
		Pathname: key,

		// file name
		Filename: key,

		// root of the path
		Root: key,
	}
}

// storage struct
type Store struct {

	// storage options
	StoreOpts
}

// A container of the paths and filename of a file
type PathKey struct {
	// the pathname to the file
	Pathname string
	// the name of the file
	Filename string
	// the root of the path
	Root string
}

// CASPathTransformFunc implements PathTransformFunc
// returns PathKey
func CASPathTransformFunc(key string) PathKey {

	// hash file key
	hash := sha1.Sum([]byte(key))

	// encode hash bytes to string
	hashStr := hex.EncodeToString(hash[:])

	// file path block sizes
	blockSize := 5

	// how many slices (folders\blocks) to file
	sliceLen := len(hashStr) / blockSize

	// path slice
	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		// create path with block size of blockSize
		// and size length of SliceLen
		from, to := i*blockSize, (i*blockSize)+blockSize

		// add block to paths slice
		paths[i] = hashStr[from:to]
	}

	// return new PathKey created from provided file key
	return PathKey{
		Pathname: strings.Join(paths, "/"), // path to file
		Filename: hashStr,                  // file name
		Root:     hashStr[:blockSize],      // file root folder
	}
}

// create new store
// return *Store
func NewStore(opts StoreOpts) *Store {

	// if custom PathTransformFunc is not provided
	if opts.PathTransformFunc == nil {
		// set PathTransformFunc to DefaultPathTransformFunc
		opts.PathTransformFunc = DefaultPathTransformFunc
	}

	// if length of storage root string is zero
	if len(opts.Root) == 0 {
		// set root to default storage root
		opts.Root = defaultRootFolder
	}

	// return store pointer
	return &Store{
		StoreOpts: opts,
	}
}

func (p *PathKey) FullPath() string {
	// format pathKey properties to full path
	return fmt.Sprintf("%s/%s", p.Pathname, p.Filename)
}

func (s *Store) Has(key string) bool {
	// transform key to PathKey
	pathKey := s.PathTransformFunc(key)

	// get file stats
	_, err := os.Stat(s.Root + "/" + pathKey.FullPath())

	// check if error returned by Stat is os.ErrNotExist
	return !errors.Is(err, os.ErrNotExist)
}

// clears all system files
func (s *Store) Clear() error {

	// remote all from storage root folder [inclusice of the root]
	return os.RemoveAll(s.Root)
}

func (s *Store) Delete(key string) error {

	// transform key to PathKey
	pathKey := s.PathTransformFunc(key)

	// defer when function ends
	defer func() {
		fmt.Printf("deleted [%s] from disk\n", pathKey.Filename)
	}()

	// delete all from file root folder
	return os.RemoveAll(s.Root + "/" + pathKey.Root)
}

// read file from storage
// return filesize (int64) | reader (io.Reader) | error
func (s *Store) Read(key string) (int64, io.Reader, error) {

	// read file stream
	n, f, err := s.readStream(key)
	if err != nil {
		// on error reading file
		return 0, nil, err
	}

	// close file on function end
	defer f.Close()

	buf := new(bytes.Buffer)

	// copy read io.Reader data to buffer
	_, err = io.Copy(buf, f)

	// return filesize (int64) | reader (io.Reader) | error
	return n, buf, err
}

func (s *Store) Write(key string, r io.Reader) (int64, error) {
	// write file stream
	return s.writeStream(key, r)
}

// read file from storage as stream
// returns file size (int64) | reader (io.Reader) | error
func (s *Store) readStream(key string) (int64, io.ReadCloser, error) {

	// transform key to PathKey
	pathKey := s.PathTransformFunc(key)

	// get file stats
	stat, err := os.Stat(s.Root + "/" + pathKey.FullPath())
	if err != nil {
		// on error getting stats
		return 0, nil, err
	}

	// open file
	r, err := os.Open(s.Root + "/" + pathKey.FullPath())

	// returns file size (int64) | reader (io.Reader) | error
	return stat.Size(), r, err
}

// write file to storage
// return written bytes size (int64) | error
func (s *Store) writeStream(key string, r io.Reader) (int64, error) {

	// transform key to PathKey
	pathKey := s.PathTransformFunc(key)

	// create directories to file
	if err := os.MkdirAll(s.Root+"/"+pathKey.Pathname, os.ModePerm); err != nil {
		return 0, err
	}

	// create file in created directories
	f, err := os.Create(s.Root + "/" + pathKey.FullPath())
	if err != nil {
		return 0, err
	}

	// write reader data to created file
	n, err := io.Copy(f, r)
	if err != nil {
		return 0, err
	}

	// return written bytes size and error
	return n, nil
}

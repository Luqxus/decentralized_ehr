package main

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathTransformFunc(t *testing.T) {
	key := "momsbestpicture"
	pathKey := CASPathTransformFunc(key)
	expected := "68044/29f74/181a6/3c50c/3d81d/733a1/2f14a/353ff"
	assert.Equal(t, expected, pathKey.Pathname)
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}

	s := NewStore(opts)
	key := "specialpicture"

	data := []byte("some jpg bytes")

	_, err := s.writeStream(key, bytes.NewReader(data))
	if err != nil {
		t.Error(err)
	}

	ok := s.Has(key)
	assert.True(t, ok)

	_, r, err := s.Read(key)
	assert.Nil(t, err)

	b, _ := io.ReadAll(r)

	assert.Equal(t, b, data)
}

func TestStoreDeleteKey(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}

	s := NewStore(opts)
	key := "specialpicture"

	data := []byte("some jpg bytes")

	_, err := s.writeStream(key, bytes.NewReader(data))
	if err != nil {
		t.Error(err)
	}

	if err := s.Delete(key); err != nil {
		t.Error(err)
	}
}

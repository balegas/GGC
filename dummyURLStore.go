package main

import (
	"bytes"
	"crypto/sha1"
	"io"
	"log"
)

type dummyURLStore struct{}

func newDummyURLStore() *dummyURLStore {
	return &dummyURLStore{}
}

func (*dummyURLStore) put(k string, v io.Reader) bool {
	if content, err := hashContent(v); err == nil {
		log.Printf("Stored: %v: %v", k, content)
		return true
	}
	log.Printf("Error: %v", k)
	return true
}

func (*dummyURLStore) get(k string) (io.Reader, error) {
	return bytes.NewReader([]byte(k)), nil
}

func hashContent(content io.Reader) (int64, error) {
	hash := sha1.New()
	return io.Copy(hash, content)
}

package main

import (
	"crypto/sha1"
	"errors"
)

var errorValueDoesNotExist = errors.New("Value does not exist")

type inMemoryURLStore struct {
	pathToHash map[string][]byte
}

func newInMemoryURLStore() *inMemoryURLStore {
	return &inMemoryURLStore{make(map[string][]byte)}
}

func (s *inMemoryURLStore) put(k string, content []byte) bool {
	s.pathToHash[k] = hashContent(content)
	return true
}

func (s *inMemoryURLStore) get(k string) ([]byte, bool) {
	if value, ok := s.pathToHash[k]; ok {
		return value, true
	}
	return nil, false
}

func hashContent(content []byte) []byte {
	hash := sha1.New()
	hash.Write(content)
	return hash.Sum(nil)
}

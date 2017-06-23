package main

import (
	"crypto/sha1"
	"errors"
	"sync"
)

var errorValueDoesNotExist = errors.New("Value does not exist")

type inMemoryURLStore struct {
	pathToHash map[string][]byte
	mutex      sync.RWMutex
}

func newInMemoryURLStore() *inMemoryURLStore {
	return &inMemoryURLStore{pathToHash: make(map[string][]byte)}
}

func (s *inMemoryURLStore) put(k string, content []byte) bool {
	h := hashContent(content)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.pathToHash[k] = h
	return true
}

func (s *inMemoryURLStore) get(k string) ([]byte, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
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

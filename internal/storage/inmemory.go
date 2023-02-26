package storage

import (
	"fmt"
	"sync"
)

var fileForSave string

type InMemory struct {
	lock sync.Mutex
	m    map[string]string
}

func NewInMemory() *InMemory {
	return &InMemory{
		m: make(map[string]string),
	}
}
func newInMemoryWithMap(innerM map[string]string) *InMemory {
	return &InMemory{
		m: innerM,
	}
}

func NewInMemoryWithFile(fileName string) *InMemory {
	if fileName == "" {
		fmt.Println("Use storage without file saving")
		return NewInMemory()
	}
	fmt.Println("Use storage with file saving")

	fileForSave = fileName

	shortURLs := restoreShortURLs(fileForSave)
	m := createMapFromShortURLs(shortURLs)
	return newInMemoryWithMap(m)

}

func (s *InMemory) GetAll() map[string]string {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.m
}

func (s *InMemory) Get(key string) (string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if v, ok := s.m[key]; ok {
		return v, nil
	}
	return "", ErrNotFound
}

func (s *InMemory) Put(key string, value string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.m[key]; ok {
		return ErrAlreadyExists
	}
	s.m[key] = value

	if saveSortURLs(fileForSave, createShortURLsFromMap(s.m)) {
		fmt.Println("ShortURLs saved in file ", fileForSave)
	}
	return nil
}

func createShortURLsFromMap(m map[string]string) *ShortURLs {
	shortURLSlice := []ShortURL{}
	for k, v := range m {
		shortURL := ShortURL{
			ID:      k,
			LongURL: v,
		}
		shortURLSlice = append(shortURLSlice, shortURL)
	}
	shortURLs := &ShortURLs{ShortURLs: shortURLSlice}
	return shortURLs
}
func createMapFromShortURLs(shortURLs *ShortURLs) map[string]string {
	m := make(map[string]string)
	if shortURLs != nil {
		shortURLsSlice := shortURLs.ShortURLs
		for i := range shortURLsSlice {
			m[shortURLsSlice[i].ID] = shortURLsSlice[i].LongURL
		}
	}
	return m
}

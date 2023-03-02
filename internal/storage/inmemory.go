package storage

import (
	"context"
	"fmt"
	"sync"
)

var fileForSave string

type InMemory struct {
	lock sync.Mutex
	m    map[string]URLRecord
}

func NewInMemory() *InMemory {
	return &InMemory{
		m: make(map[string]URLRecord),
	}
}
func newInMemoryWithMap(innerM map[string]URLRecord) *InMemory {
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

func (s *InMemory) GetAll(ctx context.Context) map[string]URLRecord {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.m
}

func (s *InMemory) GetByUser(usr string, ctx context.Context) ([]URLRecord, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	shortURLSlice := []URLRecord{}
	for _, v := range s.m {
		if usr == v.User.ID {
			shortURL := URLRecord{
				OriginalURL: v.OriginalURL,
				ShortURL:    v.ShortURL,
			}
			shortURLSlice = append(shortURLSlice, shortURL)
		}
	}
	return shortURLSlice, nil
}

func (s *InMemory) Get(key string, ctx context.Context) (URLRecord, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if v, ok := s.m[key]; ok {
		return v, nil
	}
	return URLRecord{}, ErrNotFound
}

func (s *InMemory) Put(key string, record URLRecord, ctx context.Context) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.m[key]; ok {
		return ErrAlreadyExists
	}
	s.m[key] = record

	if saveSortURLs(fileForSave, createShortURLsFromMap(s.m)) {
		fmt.Println("ShortURLs saved in file ", fileForSave)
	}
	return nil
}

func createShortURLsFromMap(m map[string]URLRecord) *URLRecords {
	shortURLSlice := []URLRecord{}
	for k, v := range m {
		shortURL := URLRecord{
			ID:          k,
			OriginalURL: v.OriginalURL,
			ShortURL:    v.ShortURL,
			User:        User{ID: v.User.ID},
		}
		shortURLSlice = append(shortURLSlice, shortURL)
	}
	shortURLs := &URLRecords{URLRecords: shortURLSlice}
	return shortURLs
}
func createMapFromShortURLs(shortURLs *URLRecords) map[string]URLRecord {
	m := make(map[string]URLRecord)
	if shortURLs != nil {
		shortURLsSlice := shortURLs.URLRecords
		for i := range shortURLsSlice {
			m[shortURLsSlice[i].ID] = shortURLsSlice[i]
		}
	}
	return m
}

func (s *InMemory) PingConnection(ctx context.Context) bool {
	return true
}
func (s *InMemory) PutAll(records []URLRecord, ctx context.Context) error {
	return nil
}
func (s *InMemory) GetByURL(url string, ctx context.Context) (URLRecord, error) {
	return URLRecord{}, nil
}

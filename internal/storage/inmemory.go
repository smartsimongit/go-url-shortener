package storage

import "sync"

type InMemory struct {
	lock sync.Mutex
	m    map[string]string
}

func NewInMemory() *InMemory {
	return &InMemory{
		m: make(map[string]string),
	}
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
	return nil
}

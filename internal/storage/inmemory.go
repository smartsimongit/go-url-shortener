package storage

import (
	"fmt"
	"sync"
)

var useFileSaving = false

type InMemory struct {
	lock sync.Mutex
	m    map[string]string
}

// NewInMemory TODO:Устарело. Удалить
func NewInMemory() *InMemory {
	return &InMemory{
		m: make(map[string]string),
	}
}

func NewInMemoryWithFile(fileName string) *InMemory {
	if fileName == "" {
		fmt.Println("Use storage without file saving")
		return NewInMemory()
	}
	fmt.Println("Use storage with file saving")
	//TODO:Восстановить из файла. Обработать что файла нет, что он пустой, что он правильно заполнен
	useFileSaving = true
	return NewInMemory()
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
	//TODO: сохранить структуру из памяти в файл, перезаписав содержимое
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.m[key]; ok {
		return ErrAlreadyExists
	}
	s.m[key] = value
	return nil
}

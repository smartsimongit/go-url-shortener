package storage

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

type Storage interface {
	Get(key string) (URLRecord, error)
	Put(key string, value URLRecord) error
	GetAll() map[string]URLRecord
	GetByUser(usr string) ([]URLRecord, error)
}

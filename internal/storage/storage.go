package storage

import (
	"context"
	"errors"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

type Storage interface {
	Get(key string, ctx context.Context) (URLRecord, error)
	Put(key string, value URLRecord, ctx context.Context) error
	GetAll(ctx context.Context) map[string]URLRecord
	GetByUser(usr string, ctx context.Context) ([]URLRecord, error)
	PingConnection(ctx context.Context) bool
	PutAll(records []URLRecord, ctx context.Context) error
	Delete(req []string, user string, ctx context.Context) error
	GetByURL(url string, ctx context.Context) (URLRecord, error)
}

type URLRecord struct {
	ID          string `json:"id,omitempty"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	User        User   `json:"usr,omitempty"`
}

type User struct {
	ID string `json:"id,omitempty"`
}

type URLRecords struct {
	URLRecords []URLRecord `json:"url_record"`
}

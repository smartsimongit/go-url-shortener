package server

import (
	"errors"
	"go-url-shortener/internal/storage"
)

var (
	ErrIncorrectPostURL     = errors.New("incorrect POST requestJSON URL")
	ErrIncorrectLongURL     = errors.New("you send incorrect LongURL")
	ErrIDParamIsMissing     = errors.New("id is missing in parameters")
	ErrIncorrectJSONRequest = errors.New("incorrect json requestJSON")
	ErrCreatedResponse      = errors.New("error created responseJson")
)

type Server struct {
	storage storage.Storage
}

func New(storage storage.Storage) *Server {
	return &Server{
		storage: storage,
	}
}

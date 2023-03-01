package server

import (
	"context"
	"github.com/gorilla/mux"
	"go-url-shortener/internal/storage"
)

type Server struct {
	ctx     context.Context
	storage storage.Storage
	repo    *storage.Repository
}

func (s *Server) AddRoutes(router *mux.Router) {
	router.HandleFunc("/{id}", s.GetHandler)
	router.HandleFunc("/", s.PostHandler)
	router.HandleFunc("/api/shorten", s.PostJSONHandler)
	router.HandleFunc("/api/user/urls", s.GetUserURLsHandler)
}

func New(ctx context.Context, storage storage.Storage) *Server {
	return &Server{
		ctx:     ctx,
		storage: storage,
	}
}
func NewWithDB(ctx context.Context, storage storage.Storage, repo *storage.Repository) *Server {
	return &Server{
		ctx:     ctx,
		storage: storage,
		repo:    repo,
	}
}

package server

import (
	"context"
	"github.com/gorilla/mux"
	"go-url-shortener/internal/storage"
)

type Server struct {
	ctx     context.Context
	storage storage.Storage
}

func (s *Server) AddRoutes(router *mux.Router) {
	router.HandleFunc("/ping", s.GetPingHandler)
	router.HandleFunc("/{id}", s.GetHandler)
	router.HandleFunc("/", s.PostHandler)
	router.HandleFunc("/api/shorten/batch", s.PostBatchURLsHandler)
	router.HandleFunc("/api/shorten", s.PostJSONHandler)
	router.HandleFunc("/api/user/urls", s.GetUserURLsHandler)

}

func New(ctx context.Context, storage storage.Storage) *Server {
	return &Server{
		ctx:     ctx,
		storage: storage,
	}
}

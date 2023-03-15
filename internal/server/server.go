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
	router.HandleFunc("/ping", s.GetPingHandler).Methods("GET")
	router.HandleFunc("/{id}", s.GetHandler).Methods("GET")
	router.HandleFunc("/", s.PostHandler).Methods("POST")
	router.HandleFunc("/api/shorten/batch", s.PostBatchURLsHandler).Methods("POST")
	router.HandleFunc("/api/shorten", s.PostJSONHandler).Methods("POST")
	router.HandleFunc("/api/user/urls", s.DeleteURLsHandler).Methods("DELETE")
	router.HandleFunc("/api/user/urls", s.GetUserURLsHandler).Methods("GET")

}

func New(ctx context.Context, storage storage.Storage) *Server {
	return &Server{
		ctx:     ctx,
		storage: storage,
	}
}

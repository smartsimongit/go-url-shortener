package main

import (
	"github.com/gorilla/mux"
	"go-url-shortener/internal/server"
	"go-url-shortener/internal/storage"
	"log"
	"net/http"
)

func main() {
	storage := storage.NewInMemory()
	router := mux.NewRouter()
	server := server.New(storage)
	router.HandleFunc("/{id}", server.GetHandler)
	router.HandleFunc("/", server.PostHandler)
	log.Fatal(http.ListenAndServe(":8080", router))
}

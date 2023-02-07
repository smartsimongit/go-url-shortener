package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"go-url-shortener/internal/server"
	"go-url-shortener/internal/storage"
)

func main() {
	store := storage.NewInMemory()
	router := mux.NewRouter()
	serv := server.New(store)
	router.HandleFunc("/{id}", serv.GetHandler)
	router.HandleFunc("/", serv.PostHandler)
	router.HandleFunc("/api/shorten", serv.PostJsonHandler)
	log.Fatal(http.ListenAndServe(":8080", router))
}

package main

import (
	"github.com/gorilla/mux"

	"log"
	"net/http"

	"go-url-shortener/internal/server"
	"go-url-shortener/internal/storage"
	"go-url-shortener/internal/util"
)

func main() {
	util.ConfigApp()
	store := storage.NewInMemoryWithFile(util.GetStorageFileName())
	router := mux.NewRouter()
	serv := server.New(store)
	router.HandleFunc("/{id}", serv.GetHandler)
	router.HandleFunc("/", serv.PostHandler)
	router.HandleFunc("/api/shorten", serv.PostJSONHandler)
	log.Fatal(http.ListenAndServe(util.GetServerAddress(), router))
}

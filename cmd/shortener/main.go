package main

import (
	"github.com/gorilla/mux"
	"go-url-shortener/internal/server"
	"log"
	"net/http"
)

var (
	repoMap map[string]string
)

func main() {
	repoMap = map[string]string{} //TODO:
	r := mux.NewRouter()
	s := server.New(repoMap)
	r.HandleFunc("/{id}", s.GetHandler)
	r.HandleFunc("/", s.PostHandler)
	log.Fatal(http.ListenAndServe(":8080", r))
}

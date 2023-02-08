package main

import (
	"github.com/gorilla/mux"
	"go-url-shortener/internal/util"
	"log"
	"net/http"

	"go-url-shortener/internal/server"
	"go-url-shortener/internal/storage"
)

func main() {
	store := storage.NewInMemory()
	router := mux.NewRouter()
	serv := server.New(store)
	router.HandleFunc("/{id}", serv.GetHandler)
	router.HandleFunc("/", serv.PostHandler)
	router.HandleFunc("/api/shorten", serv.PostJSONHandler)
	log.Fatal(http.ListenAndServe(util.GetServerAddress(), router))

	//TODO: Добавьте возможность конфигурировать сервис с помощью переменных окружения:
	//адрес запуска HTTP-сервера с помощью переменной SERVER_ADDRESS.
	//базовый адрес результирующего сокращённого URL с помощью переменной BASE_URL.
}

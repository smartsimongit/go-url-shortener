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
	store := storage.NewInMemoryWithFile(util.GetStorageFileName())
	router := mux.NewRouter()
	serv := server.New(store)
	router.HandleFunc("/{id}", serv.GetHandler)
	router.HandleFunc("/", serv.PostHandler)
	router.HandleFunc("/api/shorten", serv.PostJSONHandler)
	log.Fatal(http.ListenAndServe(util.GetServerAddress(), router))

	//TODO:Задание для трека «Сервис сокращения URL»
	//TODO:Поддержите конфигурирование сервиса с помощью флагов командной строки наравне с уже имеющимися переменными окружения:
	//TODO:флаг -a, отвечающий за адрес запуска HTTP-сервера (переменная SERVER_ADDRESS);
	//TODO:флаг -b, отвечающий за базовый адрес результирующего сокращённого URL (переменная BASE_URL);
	//TODO:флаг -f, отвечающий за путь до файла с сокращёнными URL (переменная FILE_STORAGE_PATH).
	//TODO:Во всех случаях должны быть:
	//TODO:значения по умолчанию,
	//TODO:приоритет значений, полученных через ENV, перед значениями, задаваемыми посредством флагов.

}

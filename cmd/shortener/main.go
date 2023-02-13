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
	store := storage.NewInMemoryWithFile(util.GetStorageFileName())
	router := mux.NewRouter()
	serv := server.New(store)
	router.HandleFunc("/{id}", serv.GetHandler)
	router.HandleFunc("/", serv.PostHandler)
	router.HandleFunc("/api/shorten", serv.PostJSONHandler)
	log.Fatal(http.ListenAndServe(util.GetServerAddress(), router))

	//TODO:Сохраняйте все сокращённые URL на диск в виде файла.
	//TODO:При перезапуске приложения все URL должны быть восстановлены.
	//Путь до файла должен передаваться в переменной окружения FILE_STORAGE_PATH.
	//При отсутствии переменной окружения или при её пустом значении вернитесь к хранению сокращённых URL в памяти.

	//TODO:Задание для трека «Сервис сокращения URL»
	//TODO:Поддержите конфигурирование сервиса с помощью флагов командной строки наравне с уже имеющимися переменными окружения:
	//TODO:флаг -a, отвечающий за адрес запуска HTTP-сервера (переменная SERVER_ADDRESS);
	//TODO:флаг -b, отвечающий за базовый адрес результирующего сокращённого URL (переменная BASE_URL);
	//TODO:флаг -f, отвечающий за путь до файла с сокращёнными URL (переменная FILE_STORAGE_PATH).
	//TODO:Во всех случаях должны быть:
	//TODO:значения по умолчанию,
	//TODO:приоритет значений, полученных через ENV, перед значениями, задаваемыми посредством флагов.
}

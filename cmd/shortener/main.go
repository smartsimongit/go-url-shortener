package main

import (
	"github.com/gorilla/mux"
	"go-url-shortener/internal/services"

	"log"
	"net/http"

	"go-url-shortener/internal/server"
	"go-url-shortener/internal/storage"
)

func main() {
	services.ConfigApp()
	store := storage.NewInMemoryWithFile(services.AppConfig.FileStorageURLValue)
	router := mux.NewRouter()
	serv := server.New(store)
	router.HandleFunc("/{id}", serv.GetHandler)
	router.HandleFunc("/", serv.PostHandler)
	router.HandleFunc("/api/shorten", serv.PostJSONHandler)
	log.Fatal(http.ListenAndServe(services.AppConfig.ServerAddressValue, server.Middleware(router)))

	//1.Выдавать пользователю симметрично подписанную куку, содержащую уникальный идентификатор пользователя,
	//2.если такой куки не существует или она не проходит проверку подлинности.
	//TODO: 3		Иметь хендлер GET /api/user/urls,
	//TODO: 4	 который сможет вернуть пользователю все когда-либо сокращённые им URL в формате:
	//[
	//	{
	//	"short_url": "http://...",
	//	"original_url": "http://..."
	//	},
	//	...
	//	]
	//TODO: 5	При отсутствии сокращённых пользователем URL хендлер должен отдавать HTTP-статус 204 No Content.
	//Получить куки запроса можно из поля (*http.Request).Cookie, а установить — методом http.SetCookie.

}

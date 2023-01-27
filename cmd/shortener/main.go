package main

import (
	"fmt"
	"io"
	"net/http"
)

func shortenerHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	switch r.Method {
	case "GET":
		param := r.URL.Query().Get("id")
		if param == "" {
			http.Error(w, "The id parameter is missing", http.StatusBadRequest)
			return
		}
		fmt.Println("GET METHOD with id ", param)
	case "POST":
		body, err := io.ReadAll(r.Body)
		// обрабатываем ошибку
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("POST METHOD with body ", body)
	default:
		fmt.Println(w, "Sorry, only GET and POST methods are supported.")
	}
}

func main() {
	// маршрутизация запросов обработчику
	http.HandleFunc("/", shortenerHandler)
	// запуск сервера с адресом localhost, порт 8080
	http.ListenAndServe(":8080", nil)
}

//ГОТОВ. Сервер должен быть доступен по адресу: http://localhost:8080.
//ГОТОВО. Сервер должен предоставлять два эндпоинта: POST / и GET /{id}.
//Эндпоинт POST / принимает в теле запроса строку URL для сокращения и возвращает ответ с кодом 201 и сокращённым URL в виде текстовой строки в теле.
//Эндпоинт GET /{id} принимает в качестве URL-параметра идентификатор сокращённого URL и возвращает ответ с кодом 307 и оригинальным URL в HTTP-заголовке Location.
//Нужно учесть некорректные запросы и возвращать для них ответ с кодом 400.

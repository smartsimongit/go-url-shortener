package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
)

var (
	repoMap map[string]string
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
		longUrl := repoMap[param]
		fmt.Println("mvVar is ", longUrl)
		if param == "" {
			http.Error(w, "Short url not founded", http.StatusBadRequest)
			return
		}
		w.Header().Set("Location", "longUrl")
		w.WriteHeader(http.StatusOK)

		return
	case "POST":
		body, err := io.ReadAll(r.Body)
		// обрабатываем ошибку
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println("POST METHOD with body ", string(body))
		url := string(body)
		if isUrlInvalid(url) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		genString := genString()
		repoMap[genString] = url //TODO:Проверить, может такая связка уже есть
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(genString))
		return
	default:
		http.Error(w, "Only GET and POST methods are supported", http.StatusBadRequest)
	}
}

func main() {
	repoMap = map[string]string{}

	http.HandleFunc("/", shortenerHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func isUrlInvalid(s string) bool {
	_, err := url.ParseRequestURI(s)
	if err != nil {
		return true
	}
	u, err := url.Parse(s)
	if err != nil || u.Host == "" {
		return true
	}
	return false
}

func genString() string {
	return fmt.Sprint(rand.Int63n(1000)) //TODO: переделать на строку из 6 латинских символов и цифр
}

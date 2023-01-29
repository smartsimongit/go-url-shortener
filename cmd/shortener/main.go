package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
)

var (
	repoMap map[string]string
)

func postHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		if r.URL.Path != "/" {
			http.Error(w, "url incorrect.", http.StatusBadRequest)
			return
		}
		body, err := io.ReadAll(r.Body)
		// обрабатываем ошибку
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println("POST METHOD with body ", string(body))
		url := string(body)
		if isURLInvalid(url) {
			http.Error(w, "URL is invalid", http.StatusBadRequest)
			return
		}
		genString := genString()
		repoMap[genString] = url //TODO:Проверить, может такая связка уже есть
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("http://localhost:8080/" + genString))
		return
	default:
		http.Error(w, "Only POST method is supported", http.StatusBadRequest)
	}
}
func getHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok {
			http.Error(w, "id is missing in parameters", http.StatusBadRequest)
			return
		}
		longURL := repoMap[id]
		if longURL == "" {
			http.Error(w, "Short url not founded", http.StatusBadRequest)
			return
		}
		w.Header().Set("Location", longURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	default:
		http.Error(w, "Only GET method is supported", http.StatusBadRequest)
	}
}

func main() {
	repoMap = map[string]string{}
	r := mux.NewRouter()
	r.HandleFunc("/{id}", getHandler)
	r.HandleFunc("/", postHandler)
	log.Fatal(http.ListenAndServe(":8080", r))
}

func isURLInvalid(s string) bool {
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
	symbolsForGen := "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 5)
	for i := range b {
		b[i] = symbolsForGen[rand.Intn(len(symbolsForGen))]
	}
	return string(b)

}

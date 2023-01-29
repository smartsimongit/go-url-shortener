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
		longUrl := repoMap[id]
		if longUrl == "" {
			http.Error(w, "Short url not founded", http.StatusBadRequest)
			return
		}
		w.Header().Set("Location", longUrl)
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

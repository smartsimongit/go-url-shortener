package server

import (
	"github.com/gorilla/mux"
	"go-url-shortener/internal/storage"
	"go-url-shortener/internal/util"
	"io"
	"net/http"
)

var (
	IncorrectPostUrl = "Incorrect Post request url"
	IncorrectLongURL = "You send incorrect LongURL"
	IdParamIsMissing = "Id is missing in parameters"
)

type Server struct {
	storage storage.Storage
}

func New(storage storage.Storage) *Server {
	return &Server{
		storage: storage,
	}
}

func (s *Server) PostHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		if r.URL.Path != "/" {
			http.Error(w, IncorrectPostUrl, http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		url := string(body)
		if util.IsURLInvalid(url) {
			http.Error(w, IncorrectLongURL, http.StatusBadRequest)
			return
		}

		genString := util.GenString()

		err = s.storage.Put(genString, url)
		if err == nil {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("http://localhost:8080/" + genString)) //TODO: OK
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		http.Error(w, "Only POST method is supported", http.StatusBadRequest)
	}
}

func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		//how much / in path
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok {
			http.Error(w, IdParamIsMissing, http.StatusBadRequest)
			return
		}
		longURL, err := s.storage.Get(id)
		if err == nil {
			w.Header().Set("Location", longURL)         //TODO: Header Location
			w.WriteHeader(http.StatusTemporaryRedirect) //TODO: StatusTemporaryRedirect
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return

	default:
		http.Error(w, "Only GET method is supported", http.StatusBadRequest)
	}
}

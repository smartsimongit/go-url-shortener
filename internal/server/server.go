package server

import (
	"errors"
	"io"
	"net/http"

	"github.com/gorilla/mux"

	"go-url-shortener/internal/storage"
	"go-url-shortener/internal/util"
)

var (
	ErrIncorrectPostURL = errors.New("incorrect Post request url")
	ErrIncorrectLongURL = errors.New("you send incorrect LongURL")
	ErrIDParamIsMissing = errors.New("id is missing in parameters")
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
	if r.URL.Path != "/" {
		http.Error(w, ErrIncorrectPostURL.Error(), http.StatusBadRequest)
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
		http.Error(w, ErrIncorrectLongURL.Error(), http.StatusBadRequest)
		return
	}

	genString := util.GenString()

	err = s.storage.Put(genString, url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("http://localhost:8080/" + genString))
}

func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, ErrIDParamIsMissing.Error(), http.StatusBadRequest)
		return
	}
	longURL, err := s.storage.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

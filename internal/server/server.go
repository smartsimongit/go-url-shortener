package server

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gorilla/mux"

	"go-url-shortener/internal/storage"
	"go-url-shortener/internal/util"
)

var (
	ErrIncorrectPostURL     = errors.New("incorrect Post request url")
	ErrIncorrectLongURL     = errors.New("you send incorrect LongURL")
	ErrIDParamIsMissing     = errors.New("id is missing in parameters")
	ErrIncorrectJsonRequest = errors.New("incorrect json request")
	ErrCreatedResponse      = errors.New("error created response")
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
	w.Write([]byte(util.CreateURL(genString)))
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

func (s *Server) PostJsonHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var req request
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, ErrIncorrectJsonRequest.Error(), http.StatusBadRequest)
	}

	genString := util.GenString()
	err = s.storage.Put(genString, req.Url)

	resp := response{
		ShortUrl: util.CreateURL(genString),
	}
	answer, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, ErrCreatedResponse.Error(), http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(answer)
}

type request struct {
	Url string `json:"url"`
}

type response struct {
	ShortUrl string `json:"result"`
}

//Задание для трека «Сервис сокращения URL»
//Добавьте в сервер новый эндпоинт POST /api/shorten, принимающий в теле запроса JSON-объект {"url":"<some_url>"} и возвращающий в ответ объект {"result":"<shorten_url>"}.
//Не забудьте добавить тесты на новый эндпоинт, как и на предыдущие.
//Помните про HTTP content negotiation, проставляйте правильные значения в заголовок Content-Type.

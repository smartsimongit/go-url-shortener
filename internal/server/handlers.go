package server

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"

	"go-url-shortener/internal/util"
)

func (s *Server) PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, util.ErrIncorrectPostURL.Error(), http.StatusBadRequest)
		return
	}
	encodeRequest, err := encodeBody(r)
	if err != nil {
		http.Error(w, util.ErrIncorrectJSONRequest.Error(), http.StatusBadRequest)
		return
	}
	url := string(encodeRequest)
	if util.IsURLInvalid(url) {
		http.Error(w, util.ErrIncorrectLongURL.Error(), http.StatusBadRequest)
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
		http.Error(w, util.ErrIDParamIsMissing.Error(), http.StatusBadRequest)
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

func encodeBody(r *http.Request) ([]byte, error) {
	var reader io.Reader
	if r.Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			return nil, err
		}
		reader = gz
		defer gz.Close()
	} else {
		reader = r.Body
	}
	req, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (s *Server) PostJSONHandler(w http.ResponseWriter, r *http.Request) {

	encodeRequest, err := encodeBody(r)
	if err != nil {
		http.Error(w, util.ErrIncorrectJSONRequest.Error(), http.StatusBadRequest)
		return
	}
	var req requestJSON
	err = json.Unmarshal(encodeRequest, &req)
	if err != nil {
		http.Error(w, util.ErrIncorrectJSONRequest.Error(), http.StatusBadRequest)
		return
	}

	genString := util.GenString()
	err = s.storage.Put(genString, req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp := responseJSON{
		ShortURL: util.CreateURL(genString),
	}
	answer, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, util.ErrCreatedResponse.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	w.Write(answer)
}

type requestJSON struct {
	URL string `json:"url"`
}

type responseJSON struct {
	ShortURL string `json:"result"`
}

//TODO:Вынести хэндлеры (домены?)
//TODO: Добавьте поддержку gzip
//TODO: принимать запросы в сжатом формате (HTTP-заголовок Content-Encoding);
//TODO: отдавать сжатый ответ клиенту, который поддерживает обработку сжатых ответов (HTTP-заголовок Accept-Encoding).
//TODO: Вспомните middleware из урока про HTTP-сервер, это может вам помочь.

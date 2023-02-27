package server

import (
	"encoding/hex"
	"errors"
	"fmt"
	"go-url-shortener/internal/services"

	"compress/gzip"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func (s *Server) PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, ErrIncorrectPostURL.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, ErrIncorrectJSONRequest.Error(), http.StatusBadRequest)
		return
	}
	url := string(body)
	if isURLInvalid(url) {
		http.Error(w, ErrIncorrectLongURL.Error(), http.StatusBadRequest)
		return
	}

	genString := genString()
	err = s.storage.Put(genString, url)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(createURL(genString)))
}

// TODO: Написать тест, проверить статус и хэдерц
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

func (s *Server) PostJSONHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, ErrIncorrectJSONRequest.Error(), http.StatusBadRequest)
		return
	}
	var req requestJSON
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, ErrIncorrectJSONRequest.Error(), http.StatusBadRequest)
		return
	}

	genString := genString()
	err = s.storage.Put(genString, req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp := responseJSON{
		ShortURL: createURL(genString),
	}
	answer, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, ErrCreatedResponse.Error(), http.StatusBadRequest)
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

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token, err := readCookie("Token", r)
		if err == nil && token != "" && services.VerifyToken(token) {

		} else {
			userName := genString()
			encr := services.SignName(userName)
			cookieString := hex.EncodeToString(append([]byte(userName), encr...))
			cookie := http.Cookie{Name: "token", Value: cookieString}
			http.SetCookie(w, &cookie)
			fmt.Println("encr ", encr)                 //TODO: убрать
			fmt.Println("cookieString ", cookieString) //TODO: убрать
		}

		if r.Header.Get(`Content-Encoding`) == `gzip` {
			gz, err := gzip.NewReader(r.Body)
			defer r.Body.Close()
			if err != nil {
				http.Error(w, ErrIncorrectJSONRequest.Error(), http.StatusBadRequest)
				return
			}
			r.Body = gz
			defer gz.Close()
		}

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {

			next.ServeHTTP(w, r)
			return
		}
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func genString() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	length := 8
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String() //
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

func createURL(s string) string {
	return services.AppConfig.BaseURLValue + "/" + s
}

func readCookie(name string, r *http.Request) (value string, err error) {
	if name == "" {
		return value, errors.New("you are trying to read empty cookie")
	}
	cookie, err := r.Cookie(name)
	if err != nil {
		return value, err
	}
	str := cookie.Value
	value, _ = url.QueryUnescape(str)
	return value, err
}

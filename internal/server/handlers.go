package server

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"compress/gzip"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go-url-shortener/internal/services"
	"go-url-shortener/internal/storage"
)

func (s *Server) GetPingHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	if !s.storage.PingConnection(ctx) {
		http.Error(w, ErrPingConnection.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) GetUserURLsHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	user, err := getUser(r)

	if err != nil || user == "" {
		http.Error(w, ErrServer.Error(), http.StatusInternalServerError)
		return
	}
	records, err := s.storage.GetByUser(user, ctx)

	if err != nil || len(records) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	urlPair := recordsToURLDto(records)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	answer, err := json.Marshal(urlPair)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write(answer)
}

func recordsToURLDto(records []storage.URLRecord) []URLDto {
	var uRLDtoSlice []URLDto
	for _, v := range records {
		uRLDto := URLDto{
			OriginalURL: v.OriginalURL,
			ShortURL:    v.ShortURL,
		}
		uRLDtoSlice = append(uRLDtoSlice, uRLDto)
	}
	return uRLDtoSlice
}

type URLDto struct {
	ShortURL    string `json:"short_url,omitempty"`
	OriginalURL string `json:"original_url,omitempty"`
}

func (s *Server) PostHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	if r.URL.Path != "/" {
		http.Error(w, ErrIncorrectPostURL.Error(), http.StatusBadRequest)
		return
	}
	user, err := getUser(r)
	if err != nil && user != "" {
		http.Error(w, ErrIncorrectJSONRequest.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, ErrIncorrectJSONRequest.Error(), http.StatusBadRequest)
		return
	}
	incomingURL := string(body)
	if isURLInvalid(incomingURL) {
		http.Error(w, ErrIncorrectLongURL.Error(), http.StatusBadRequest)
		return
	}
	genString := genString()
	rec := storage.URLRecord{ID: genString,
		ShortURL:    createURL(genString),
		OriginalURL: incomingURL,
		User:        storage.User{ID: user}}
	httpStatus := http.StatusCreated
	recBd, err := s.storage.GetByURL(incomingURL, ctx)
	if err != nil {
		err = s.storage.Put(genString, rec, ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(httpStatus)
		w.Write([]byte(rec.ShortURL))
		return
	}
	rec = recBd
	httpStatus = http.StatusConflict
	w.WriteHeader(httpStatus)
	w.Write([]byte(rec.ShortURL))
}

func getUser(r *http.Request) (string, error) {
	token, err := readCookie("token", r)
	if err != nil {
		return "", err
	}
	user := services.GetUserFromToken(token)
	return user, nil
}

func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, ErrIDParamIsMissing.Error(), http.StatusBadRequest)
		return
	}
	longURL, err := s.storage.Get(id, ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", longURL.OriginalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (s *Server) PostJSONHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	user, err := getUser(r)
	if err != nil && user != "" {
		http.Error(w, ErrIncorrectJSONRequest.Error(), http.StatusBadRequest)
		return
	}
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
	rec := storage.URLRecord{ID: genString,
		ShortURL:    createURL(genString),
		OriginalURL: req.URL,
		User:        storage.User{ID: user}}
	httpStatus := http.StatusCreated
	recBd, err := s.storage.GetByURL(req.URL, ctx)
	if err == nil {
		rec = recBd
		httpStatus = http.StatusConflict
	} else {
		err = s.storage.Put(genString, rec, ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	resp := responseJSON{
		ShortURL: rec.ShortURL,
	}
	answer, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, ErrCreatedResponse.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)

	w.Write(answer)
}

type requestJSON struct {
	URL string `json:"url"`
}

type responseJSON struct {
	ShortURL string `json:"result"`
}

func (s *Server) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := getTokenCookie(r)
		if err != nil || token == "" || !services.VerifyToken(token) {
			userName := genString()
			log.Info().Str("new_user", userName)

			encr := services.SignName(userName)
			cookieString := hex.EncodeToString(append([]byte(userName), encr...))
			cookie := http.Cookie{Name: services.AuthnCookieName, Value: cookieString}
			http.SetCookie(w, &cookie)
			r.AddCookie(&cookie)
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
func getTokenCookie(r *http.Request) (value string, err error) {
	return readCookie(services.AuthnCookieName, r)
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

type urlBatchReq struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}
type urlBatchResp struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func (s *Server) PostBatchURLsHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	user, err := getUser(r)
	if err != nil && user != "" {
		http.Error(w, "User is expected", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, ErrIncorrectJSONRequest.Error(), http.StatusBadRequest)
		return
	}
	var req []urlBatchReq
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, ErrIncorrectJSONRequest.Error(), http.StatusBadRequest)
		return
	}
	var records []storage.URLRecord
	var urlBatchResponses []urlBatchResp

	for _, url := range req {

		genString := genString()
		shortURL := createURL(genString)
		urlBatchResp := urlBatchResp{
			CorrelationID: url.CorrelationID,
			ShortURL:      shortURL,
		}
		urlBatchResponses = append(urlBatchResponses, urlBatchResp)
		rec := storage.URLRecord{
			ID:          genString,
			ShortURL:    shortURL,
			OriginalURL: url.OriginalURL,
			User:        storage.User{ID: user},
		}
		records = append(records, rec)
	}

	err = s.storage.PutAll(records, ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	answer, err := json.Marshal(urlBatchResponses)
	if err != nil {
		http.Error(w, ErrCreatedResponse.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	w.Write(answer)
}

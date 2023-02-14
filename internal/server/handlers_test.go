package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go-url-shortener/internal/storage"
	"go-url-shortener/internal/util"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (int, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	defer resp.Body.Close()
	return resp.StatusCode, string(respBody)
}

func testPOSTResponse(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) *http.Response {
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}

func TestHandlers_PostHandlerOk(t *testing.T) {
	util.ConfigApp()

	sendedURL := "https://practicum.yandex.ru/"
	expectedStatus := http.StatusCreated
	path := "/"

	server := New(storage.NewInMemory())
	r := mux.NewRouter()
	ts := httptest.NewServer(r)
	r.HandleFunc("/", server.PostHandler)
	defer ts.Close()

	statusCode, body := testRequest(t, ts, "POST", path, bytes.NewBuffer([]byte((sendedURL))))
	assert.Equal(t, expectedStatus, statusCode)
	assert.NotEmpty(t, body)
	fmt.Println("body is ", body)
	assert.False(t, util.IsURLInvalid(body))

}
func TestHandlers_PostHandlerErrorStatus(t *testing.T) {
	type want struct {
		code     int
		response string
	}
	tests := []struct {
		name     string
		reqID    string
		longLink string
		target   string
		want     want
	}{
		{
			name:     "test #1 You send incorrect LongURL",
			longLink: "JsdjjsSJdsS",
			target:   "/",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:     "test #3 page not found",
			longLink: "http://mail.yandex.ru",
			target:   "/target/",
			want: want{
				code: http.StatusNotFound,
			},
		},
	}
	server := New(storage.NewInMemory())
	r := mux.NewRouter()
	ts := httptest.NewServer(r)
	r.HandleFunc("/", server.PostHandler)
	r.HandleFunc("/{id}", server.GetHandler)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode, body := testRequest(t, ts, "POST", tt.target, bytes.NewBuffer([]byte((tt.longLink))))
			assert.Equal(t, tt.want.code, statusCode)
			assert.NotEmpty(t, body)
		})
	}
}

func TestHandlers_GetHandlerErrorStatus(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name   string
		reqID  string
		target string
		want   want
	}{
		{
			name:   "test #1  not found code is storage",
			target: "/sd3rt",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:   "test #2 page not found",
			target: "/dsgdsfsd/",
			want: want{
				code: http.StatusNotFound,
			},
		},
		{
			name:   "test #3 Only POST method for this url",
			target: "/",
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}

	server := New(storage.NewInMemory())
	r := mux.NewRouter()
	ts := httptest.NewServer(r)
	r.HandleFunc("/", server.PostHandler)
	r.HandleFunc("/{id}", server.GetHandler)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode, body := testRequest(t, ts, "GET", tt.target, nil)
			assert.Equal(t, tt.want.code, statusCode)
			assert.NotEmpty(t, body)
		})
	}
}

func TestHandlers_PostJSONHandlerErrorStatus(t *testing.T) {
	type want struct {
		code     int
		response string
	}
	tests := []struct {
		name string
		body string
		want want
	}{
		{
			name: "test #1 incorrect json request",
			body: "",
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}
	path := "/api/shorten"
	method := "POST"

	server := New(storage.NewInMemory())
	r := mux.NewRouter()
	ts := httptest.NewServer(r)
	r.HandleFunc(path, server.PostJSONHandler)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode, body := testRequest(t, ts, method, path, bytes.NewBuffer([]byte((tt.body))))
			fmt.Println(body)
			assert.Equal(t, tt.want.code, statusCode)
			assert.NotEmpty(t, body)
		})
	}
}

func TestHandlers_PostJSONHandlerOKStatus(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name    string
		bodyReq string
		want    want
	}{
		{
			name:    "test #1 ok response",
			bodyReq: "{\"url\":\"https://practicum.yandex.ru/\"}",
			want: want{
				code:        http.StatusCreated,
				contentType: "application/json",
			},
		},
	}
	path := "/api/shorten"
	method := "POST"

	server := New(storage.NewInMemory())
	r := mux.NewRouter()
	ts := httptest.NewServer(r)
	r.HandleFunc(path, server.PostJSONHandler)
	defer ts.Close()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := testPOSTResponse(t, ts, method, path, bytes.NewBuffer([]byte((tt.bodyReq))))
			respBody, err := io.ReadAll(resp.Body)
			assert.Nil(t, err)
			defer resp.Body.Close()
			contentType := resp.Header.Get("Content-Type")
			fmt.Println(string(respBody))
			assert.NotEmpty(t, respBody)
			assert.True(t, isJSON(string(respBody)))
			assert.Equal(t, tt.want.code, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, contentType)
		})
	}

}
func isJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

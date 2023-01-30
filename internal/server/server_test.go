package server

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go-url-shortener/internal/storage"
	"go-url-shortener/internal/util"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
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

func TestServer_PostHandlerOk(t *testing.T) {

	sendedURL := "https://practicum.yandex.ru/"
	expectedStatus := http.StatusCreated
	path := "/"

	server := New(storage.NewInMemory())
	r := mux.NewRouter()
	ts := httptest.NewServer(r)
	r.HandleFunc("/", server.PostHandler)
	defer ts.Close()

	statusCode, body := testRequest(t, ts, "POST", path, bytes.NewBuffer([]byte((sendedURL))))
	fmt.Println("body is ", body)
	assert.Equal(t, expectedStatus, statusCode)
	assert.NotEmpty(t, body)
	assert.False(t, util.IsURLInvalid(body))

}
func TestServer_PostHandlerErrorStatus(t *testing.T) {
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
			fmt.Println("body is ", body)
			assert.Equal(t, tt.want.code, statusCode)
			assert.NotEmpty(t, body)
		})
	}
}

func TestServer_GetHandlerErrorStatus(t *testing.T) {
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
			fmt.Println("body is ", body)
			assert.Equal(t, tt.want.code, statusCode)
			assert.NotEmpty(t, body)
		})
	}
}

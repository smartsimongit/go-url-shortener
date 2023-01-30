package server

import (
	"bytes"
	"fmt"
	"go-url-shortener/internal/storage"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TODO: Дописать тесты
func TestServer_PostHandler(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name     string
		longLink string
		target   string
		server   *Server
		want     want
	}{
		{
			name:     "test #1",
			longLink: "https://practicum.yandex.ru/",
			target:   "/",
			server:   New(storage.NewInMemory()),
			want: want{
				code: 201,
			},
		},
		{
			name:     "test #2",
			longLink: "JsdjjsSJdsS",
			target:   "/",
			server:   New(storage.NewInMemory()),
			want: want{
				code: 400,
			},
		},
		{
			name:     "test #3",
			longLink: "https://practicum.yandex.ru/",
			target:   "/target/",
			server:   New(storage.NewInMemory()),
			want: want{
				code: 400,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.target, bytes.NewBuffer([]byte(tt.longLink)))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(tt.server.PostHandler)
			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()
			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}
			body, _ := io.ReadAll(res.Body)
			fmt.Println("Req body is ", body)
			if body == nil {
				t.Errorf("Incorrect response body")
			}
		})
	}
}

func TestServer_GetHandler(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name     string
		reqID    string
		longLink string
		target   string
		server   *Server
		want     want
	}{
		{
			name:     "test #1",
			longLink: "https://practicum.yandex.ru/",
			target:   "/",
			server:   New(storage.NewInMemoryWithInnerMap(map[string]string{"sdsSs": "https://practicum.yandex.ru/"})),
			want: want{
				code: http.StatusBadRequest,
			},
		},
		//{
		//	name:     "test #2",
		//	longLink: "JsdjjsSJdsS",
		//	target:   "/",
		//	server:   New(storage.NewInMemory()),
		//	want: want{
		//		code: 400,
		//	},
		//},
		//{
		//	name:     "test #3",
		//	longLink: "https://practicum.yandex.ru/",
		//	target:   "/target/",
		//	server:   New(storage.NewInMemory()),
		//	want: want{
		//		code: 400,
		//	},
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.target, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(tt.server.GetHandler)
			h.ServeHTTP(w, request)

			res := w.Result()
			defer res.Body.Close()
			body, _ := io.ReadAll(res.Body)
			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}
			fmt.Println("Req body is ", body)
			if body == nil {
				t.Errorf("Incorrect response body")
			}
		})
	}
}

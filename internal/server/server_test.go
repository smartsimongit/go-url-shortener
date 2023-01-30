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

			// создаём новый Recorder
			w := httptest.NewRecorder()
			// определяем хендлер
			h := http.HandlerFunc(tt.server.PostHandler)
			// запускаем сервер
			h.ServeHTTP(w, request)

			res := w.Result()
			defer res.Body.Close()
			// проверяем код ответа
			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}
			body, _ := io.ReadAll(res.Body)
			fmt.Println(string(body))
		})
	}
}

func TestServer_GetHandlerHandler(t *testing.T) {

}

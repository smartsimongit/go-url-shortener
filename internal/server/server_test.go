package server

import (
	"go-url-shortener/internal/storage"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_getHandler(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name   string
		reqID  string
		server Server
		want   want
	}{
		{
			name:  "negative test #1",
			reqID: "sdf3p",
			server: Server{
				storage: storage.NewInMemory(),
			},
			want: want{
				code: 400,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/"+tt.reqID, nil)

			// создаём новый Recorder
			w := httptest.NewRecorder()
			// определяем хендлер
			h := http.HandlerFunc(tt.server.GetHandler)
			// запускаем сервер
			h.ServeHTTP(w, request)

			res := w.Result()
			defer res.Body.Close()
			// проверяем код ответа
			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}
		})
	}
}

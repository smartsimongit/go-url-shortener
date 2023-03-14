package main

import (
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go-url-shortener/internal/services"

	"context"
	"net/http"

	"go-url-shortener/internal/server"
	"go-url-shortener/internal/storage"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	ctx := context.Background()
	services.ConfigApp()
	store := storage.NewInMemory(storage.WithFile(services.AppConfig.FileStorageURLValue))
	router := mux.NewRouter()

	serv := server.New(ctx, store)
	dbpool, err := storage.New(ctx)
	if err == nil {
		log.Info().Msg("init server with DB")
		serv = server.New(ctx, storage.NewRepository(dbpool))
	}
	serv.AddRoutes(router)
	log.Fatal().Err(http.ListenAndServe(services.AppConfig.ServerAddressValue, serv.Middleware(router)))

	//	TODO: Фактический результат удаления может происходить позже
	//	TODO: Используйте паттерн fanIn для максимального наполнения буфера объектов обновления.
}

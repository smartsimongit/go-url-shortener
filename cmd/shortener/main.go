package main

import (
	"context"
	"github.com/gorilla/mux"
	"go-url-shortener/internal/services"

	"log"
	"net/http"

	"go-url-shortener/internal/server"
	"go-url-shortener/internal/storage"
)

func main() {
	ctx := context.Background()

	dbpool, err := storage.InitDBConn(ctx)
	if err != nil {
		log.Fatalf("%w failed to init DB connection", err)
	}

	services.ConfigApp()
	store := storage.NewInMemoryWithFile(services.AppConfig.FileStorageURLValue)
	router := mux.NewRouter()
	serv := server.NewWithDB(ctx, store, storage.NewRepository(dbpool))

	serv.AddRoutes(router)
	log.Fatal(http.ListenAndServe(services.AppConfig.ServerAddressValue, serv.Middleware(router)))

}

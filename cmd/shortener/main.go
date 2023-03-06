package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"go-url-shortener/internal/services"

	"log"
	"net/http"

	"go-url-shortener/internal/server"
	"go-url-shortener/internal/storage"
)

func main() {
	ctx := context.Background()
	services.ConfigApp()
	store := storage.NewInMemoryWithFile(services.AppConfig.FileStorageURLValue)
	router := mux.NewRouter()

	serv := server.New(ctx, store)
	dbpool, err := storage.New(ctx)
	if err == nil {
		fmt.Println("init server with BD")
		serv = server.New(ctx, storage.NewRepository(dbpool))
	}
	serv.AddRoutes(router)
	log.Fatal(http.ListenAndServe(services.AppConfig.ServerAddressValue, serv.Middleware(router)))

}

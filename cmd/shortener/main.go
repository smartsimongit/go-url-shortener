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

	dbpool, err := storage.InitDBConn(ctx)
	serv := server.New(ctx, store)
	if err == nil {
		fmt.Println("init server with BD")
		serv = server.NewWithDB(ctx, store, storage.NewRepository(dbpool))
	}
	serv.AddRoutes(router)
	log.Fatal(http.ListenAndServe(services.AppConfig.ServerAddressValue, serv.Middleware(router)))

	//Перепишите сервис так, чтобы СУБД PostgreSQL стала хранилищем сокращённых URL вместо текущей реализации.
	//	Сервису нужно самостоятельно создать все необходимые таблицы в базе данных. Схема и формат хранения остаются на ваше усмотрение.
	//	При отсутствии переменной окружения DATABASE_DSN или флага командной строки -d или при их пустых значениях вернитесь последовательно к:
	//хранению сокращённых URL в файле при наличии соответствующей переменной окружения или флага командной строки;
	//хранению сокращённых URL в памяти.

}

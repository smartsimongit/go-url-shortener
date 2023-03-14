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

	//	Задание для трека «Сервис сокращения URL»

	//	 1Сделайте в таблице базы данных с сокращёнными URL дополнительное поле с флагом, указывающим на то,
	//	 1что URL должен считаться удалённым.

	//	 2Далее добавьте в сервис новый асинхронный хендлер DELETE /api/user/urls,
	//	 2который принимает список идентификаторов сокращённых URL для удаления в формате:
	//   2[ "a", "b", "c", "d", ...]

	//	3В случае успешного приёма запроса хендлер должен возвращать HTTP-статус 202 Accepted.

	//	TODO:4Фактический результат удаления может происходить позже
	//	TODO:4— каким-либо образом оповещать пользователя об успешности или неуспешности не нужно.

	//	TODO:5Успешно удалить URL может пользователь, его создавший.

	//	TODO:5При запросе удалённого URL с помощью хендлера GET /{id} нужно вернуть статус 410 Gone.

	//	Совет:
	//	Для эффективного проставления флага удаления в базе данных используйте множественное обновление (batch update).
	//	Используйте паттерн fanIn для максимального наполнения буфера объектов обновления.
}

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
	//	Сделайте в таблице базы данных с сокращёнными URL дополнительное поле с флагом, указывающим на то, что URL должен считаться удалённым.
	//		Далее добавьте в сервис новый асинхронный хендлер DELETE /api/user/urls, который принимает список идентификаторов сокращённых URL для удаления в формате:
	//[ "a", "b", "c", "d", ...]
	//	В случае успешного приёма запроса хендлер должен возвращать HTTP-статус 202 Accepted.
	//	Фактический результат удаления может происходить позже — каким-либо образом оповещать пользователя об успешности или неуспешности не нужно.
	//	Успешно удалить URL может пользователь, его создавший. При запросе удалённого URL с помощью хендлера GET /{id} нужно вернуть статус 410 Gone.
	//	Совет:
	//	Для эффективного проставления флага удаления в базе данных используйте множественное обновление (batch update).
	//	Используйте паттерн fanIn для максимального наполнения буфера объектов обновления.
}

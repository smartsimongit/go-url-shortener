package storage

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"go-url-shortener/internal/services"
	"log"

	"context"
	"fmt"
	"net"
	"time"
)

func InitDBConn(ctx context.Context) (dbpool *pgxpool.Pool, err error) {

	url := services.AppConfig.DBAddressURL
	//url := "postgres://postgres:postgres@localhost:5432/postgres"

	if url == "" {
		err = fmt.Errorf("failed to get url: %w", err)
		return
	}

	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		err = fmt.Errorf("failed to parse pg config: %w", err)
		return
	}
	cfg.MaxConns = int32(5)
	cfg.MinConns = int32(1)
	cfg.HealthCheckPeriod = 1 * time.Minute
	cfg.MaxConnLifetime = 24 * time.Hour
	cfg.MaxConnIdleTime = 30 * time.Minute
	cfg.ConnConfig.ConnectTimeout = 1 * time.Second
	cfg.ConnConfig.DialFunc = (&net.Dialer{
		KeepAlive: cfg.HealthCheckPeriod,
		Timeout:   cfg.ConnConfig.ConnectTimeout,
	}).DialContext

	dbpool, err = pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		err = fmt.Errorf("failed to connect config: %w", err)
		return
	}

	return

}

// REPO
type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	r := &Repository{pool: pool}
	r.createTables()
	return r
}

func (r *Repository) PingConnection(ctx context.Context) bool {
	err := r.pool.Ping(ctx)
	return err == nil
}

func (r *Repository) createTables() {
	ctx := context.Background()
	_, err := r.pool.Exec(ctx, "create table if not exists public.link_pairs(id varchar(64) primary key, short_url    varchar(64)  not null, original_url varchar(256) not null, usr varchar(64)  not null);")
	if err != nil {
		fmt.Println("таблица не создалась ", err.Error())
		log.Fatal(err)
	}
}
func (r *Repository) Get(key string, ctx context.Context) (URLRecord, error) {
	record := URLRecord{}

	row := r.pool.QueryRow(ctx,
		"SELECT lp.id, lp.short_url, lp.original_url, lp.usr FROM public.link_pairs lp WHERE lp.id = $1",
		key)

	err := row.Scan(&record.ID, &record.ShortURL, &record.OriginalURL, &record.User.ID)
	if err != nil {
		fmt.Println("ошибка чтения ", err.Error())
		return record, err
	}
	return record, err
}

func (r *Repository) Put(key string, value URLRecord, ctx context.Context) error {
	_, err := r.pool.Exec(ctx,
		"INSERT INTO public.link_pairs (id, short_url, original_url, usr) VALUES($1,$2,$3,$4)",
		key, value.ShortURL, value.OriginalURL, value.User.ID)
	if err != nil {
		fmt.Println("ошибка записи ", err.Error())
		return err
	}
	return nil
}
func (r *Repository) GetAll(ctx context.Context) map[string]URLRecord {
	return nil //TODO:
}
func (r *Repository) GetByUser(usr string, ctx context.Context) ([]URLRecord, error) {
	shortURLSlice := []URLRecord{}
	rows, err := r.pool.Query(ctx,
		"SELECT lp.id, lp.short_url, lp.original_url, lp.usr FROM public.link_pairs lp WHERE lp.usr = $1",
		usr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	defer rows.Close()

	// пробегаем по всем записям
	for rows.Next() {
		var r URLRecord
		err = rows.Scan(&r.ID, &r.ShortURL, &r.OriginalURL, &r.User.ID)
		if err != nil {
			return nil, err
		}
		shortURLSlice = append(shortURLSlice, r)
	}
	return shortURLSlice, nil
}

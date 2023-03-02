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

	url := services.AppConfig.DBAddressURL //url := "postgres://postgres:postgres@localhost:5432/postgres"

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
	_, err := r.pool.Query(ctx, "create table if not exists public.link_pairs(id varchar(64) primary key, short_url    varchar(64)  not null, original_url varchar(256) not null, usr varchar(64)  not null);")
	if err != nil {
		fmt.Println("таблица не создалась ", err.Error())
		log.Fatal(err)
	}
}

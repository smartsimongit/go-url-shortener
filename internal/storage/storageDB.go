package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"go-url-shortener/internal/services"
	"net"
	"time"
)

func InitDBConn(ctx context.Context) (dbpool *pgxpool.Pool, err error) {

	url := services.AppConfig.DBAddressURL

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
	return &Repository{pool: pool}
}

func (r *Repository) PingConnection(ctx context.Context) bool {
	err := r.pool.Ping(ctx)
	return err == nil
}

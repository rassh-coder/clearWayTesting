package repository

import (
	"clearWayTest/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

func NewPostgresDB(cfg *config.Config) (*pgx.Conn, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.PG.Username, cfg.PG.Password, cfg.PG.Host, cfg.PG.Port, cfg.PG.DBName)
	conn, err := pgx.Connect(context.Background(), connStr)

	if err != nil {
		return nil, err
	}

	return conn, nil
}

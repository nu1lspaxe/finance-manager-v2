package server

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SetPGConn(ctx context.Context) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, os.Getenv("POSTGRES_URL"))
	if err != nil {
		return nil, err
	}

	if connErr := pool.Ping(ctx); connErr != nil {
		return nil, connErr
	}

	return pool, nil
}

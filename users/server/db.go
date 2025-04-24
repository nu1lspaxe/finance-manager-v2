package server

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func GetPgPool(ctx context.Context, pgLink string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, pgLink)
	if err != nil {
		return nil, err
	}

	if connErr := pool.Ping(ctx); connErr != nil {
		return nil, connErr
	}

	return pool, nil
}

func GetRedisClient(ctx context.Context, redisLink string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisLink,
		Password: "",
		DB:       0,
	})

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return client, nil
}

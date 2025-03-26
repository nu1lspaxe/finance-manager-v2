package pkg

import (
	"context"
	"time"
)

func (u *UserService) redisHSet(ctx context.Context, key string, expiration time.Duration, values ...any) error {
	err := u.redisClient.HSet(ctx, key, values).Err()
	if err != nil {
		return err
	}

	err = u.redisClient.Expire(ctx, key, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) redisHDel(ctx context.Context, key string, field string) error {
	err := u.redisClient.HDel(ctx, key, field).Err()
	if err != nil {
		return err
	}
	return nil
}

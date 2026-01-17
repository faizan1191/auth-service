package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

func SetRefreshToken(rdb *redis.Client, token, userID string) error {
	return rdb.Set(ctx, token, userID, 7*24*time.Hour).Err()
}

func GetUserIDByRefreshToken(rdb *redis.Client, token string) (string, error) {
	return rdb.Get(ctx, token).Result()
}

func DeleteRefreshToken(rdb *redis.Client, token string) error {
	return rdb.Del(ctx, token).Err()
}

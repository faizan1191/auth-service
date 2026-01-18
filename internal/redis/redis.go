package redis

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
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

func SetResetToken(rdb *redis.Client, token, userID string) error {
	hashedToken := hashToken(token)
	key := "reset" + hashedToken
	return rdb.Set(ctx, key, userID, 15*time.Minute).Err()
}

func ConsumeResetToken(rdb *redis.Client, token string) (string, error) {
	hashedToken := hashToken(token)
	key := "reset" + hashedToken

	userID, err := rdb.Get(ctx, key).Result()

	if err == redis.Nil {
		return "", nil // invalid or expired
	}
	if err != nil {
		return "", err
	}

	// one-time use
	rdb.Del(ctx, key)

	return userID, nil
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

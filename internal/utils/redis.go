package utils

import (
	"context"
	"log"
	"time"

	"github.com/aakritigkmit/payment-gateway/internal/config"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() {
	cfg := config.GetConfig()

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	if err := PingRedis(context.Background()); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	} else {
		log.Println("Redis connected successfully")
	}
}

func GetRedisClient() *redis.Client {
	return RedisClient
}

func PingRedis(ctx context.Context) error {
	return RedisClient.Ping(ctx).Err()
}

// SetRedisKey sets a value in Redis with an optional expiration
func SetRedisKey(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return RedisClient.Set(ctx, key, value, expiration).Err()
}

// GetRedisKey retrieves a value from Redis by key
func GetRedisKey(ctx context.Context, key string) (string, error) {
	return RedisClient.Get(ctx, key).Result()
}

// DeleteRedisKey deletes a key from Redis
func DeleteRedisKey(ctx context.Context, key string) error {
	return RedisClient.Del(ctx, key).Err()
}

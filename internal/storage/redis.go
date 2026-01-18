package storage

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

// NewRedisClient initializes the Redis client and pings the server to ensure connectivity.
func NewRedisClient(addr string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // Default Docker Redis has no password
		DB:       0,  // Default DB
	})

	// Pinging Redis to confirm connection before proceeding
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis at %s: %v", addr, err)
	}

	return rdb
}

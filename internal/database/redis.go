package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nighbee/evently/internal/config"
	"github.com/redis/go-redis/v9"
)

// ConnectRedis creates a Redis client connection
func ConnectRedis(cfg config.Config) *redis.Client {
	addr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.RedisPassword,
		DB:       0, // use default DB
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("warning: failed to connect to Redis: %v", err)
		return nil
	}

	log.Printf("connected to Redis at %s", addr)
	return rdb
}

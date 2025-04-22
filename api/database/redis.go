package database

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
)

func InitRedis(redisUrl string) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: redisUrl,
	})

	if err := RedisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	fmt.Println("[Redis] Connected successfully 🚀")
}

func CloseRedis() {
	if err := RedisClient.Close(); err != nil {
		log.Printf("Error closing Redis: %v", err)
	}
	fmt.Println("[Redis] Connection closed gracefully 🧹")
}

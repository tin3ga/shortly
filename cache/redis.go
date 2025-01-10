package cache

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"

	"github.com/tin3ga/shortly/utils"
)

func InitializeRedis(ctx context.Context, redisAddr string, redisPassword string, redisDB string) (*redis.Client, error) {
	// Convert Redis DB from string to int
	redisDBInt, err := utils.ConvertStr(redisDB)
	if err != nil {
		log.Printf("redisDB %v", err)
	}
	// Create Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDBInt,
	})

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Printf("Cannot connect to redis db: \nError %v", err)
		return nil, err
	}

	if pong == "PONG" {
		log.Print("Connected to cache server")

	}

	return rdb, nil
}

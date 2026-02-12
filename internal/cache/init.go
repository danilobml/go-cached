package cache

import (
	"context"
	"log"
	"time"
)

func InitRedis(cfg Config) (*RedisClient, error) {
	client := NewRedisClient(cfg.Addr, cfg.Password, cfg.DB)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := client.Ping(ctx); err != nil {
		return nil, err
	}

	log.Println("Connected to Redis!")

	return client, nil
}

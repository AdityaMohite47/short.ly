package db

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func Connectredis(ctx context.Context, url string) (*redis.Client, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("invalid redis url: %w", err)
	}

	rdb := redis.NewClient(opts)

	// reachability check
	if err := rdb.Ping(ctx).Err(); err != nil {
		rdb.Close()
		return nil, fmt.Errorf("error pinging redis: %w", err)
	}

	fmt.Println("Connected to Redis")
	return rdb, nil
}

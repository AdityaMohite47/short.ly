package db

import (
	"context"
	"fmt"
	"time"
	"os"
	"strconv"
	"log"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func CacheTTL() (time.Duration, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("Failed Loading .env variables", err)
	}
	s := os.Getenv("CACHE_URL_EXPIRY_HOURS")
	if s == "" {
		return 24 * time.Hour, nil
	}

	hours, err := strconv.Atoi(s)
	if err != nil || hours <= 0 {
		return 0, fmt.Errorf("CACHE_URL_EXPIRY_HOURS must be a positive whole number, got %q", s)
	}
	return time.Duration(hours) * time.Hour, nil
}

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
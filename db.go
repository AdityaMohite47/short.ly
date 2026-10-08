package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type DBconn struct {
	Pool *pgxpool.Pool
}

func (conn *DBconn) ConnectDB() {
	if err := godotenv.Load(); err != nil {
		log.Println("Failed Loading .env variables", err)
		return
	}
	url := os.Getenv("DB_URL")

	// defining a connection pool
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		log.Println("Couldn't create Pool", err)
		return
	}

	// DB reachability check
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		log.Println("Error pining the Database", err, "/n Connection closed...")
		return
	}

	conn.Pool = pool
	fmt.Println("Connected to PostgreSQL")
}

var ErrNotFound = errors.New("link not found")

// Save a long URL
func (d *DBconn) getShortURL(ctx context.Context, longURL string) (string, error) {
	// saving long URL
	var id int64
	if err := d.Pool.QueryRow(ctx,
		"INSERT INTO short_links (long_urls) VALUES ($1) RETURNING id",
		longURL,
	).Scan(&id); err != nil {
		return "", err
	}

	short_url := os.Getenv("DomainName") + Encode(id)
	fmt.Println(short_url)

	if _, err := d.Pool.Exec(ctx, "INSERT INTO short_links (short_urls) VALUES ($1)", short_url); err != nil {
		return "", err
	}
	return short_url, nil

}

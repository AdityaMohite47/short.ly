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
		log.Println("DB URL not set", err)
		return
	}
	url := os.Getenv("DB_URL")

	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		log.Println("Couldn't create Pool", err)
		return
	}

	// Actually check the database is reachable
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
func (d *DBconn) InsertLink(ctx context.Context, longURL string) error {
	_, err := d.Pool.Exec(ctx,
		`INSERT INTO short_links (long_urls) VALUES ($1)`,
		longURL,
	)
	return err
}

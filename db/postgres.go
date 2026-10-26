package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"short.ly/shortcode"
)


// Connect opens a connection pool to url and checks the database is reachable
func Connectpsql(ctx context.Context, url string) (*pgxpool.Pool, error) {
	// defining a connection pool
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("couldn't create pool: %w", err)
	}

	// DB reachability check
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("error pinging the database: %w", err)
	}

	fmt.Println("Connected to PostgreSQL")
	return  pool, nil
}


var ErrNotFound = errors.New("link not found")

// CreateLink saves a long URL under a freshly generated random code and returns the code
func CreateLink(pool *pgxpool.Pool , ctx context.Context, longURL string) (string, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("Failed Loading .env variables", err)
	}
	for range os.Getenv("CODE_MAX_RETRIES") {
		code, err := shortcode.GenerateCode(shortcode.CodeLength)
		if err != nil {
			return "", err
		}

		_, err = pool.Exec(ctx,
			"INSERT INTO shortlinks (code, long_url) VALUES ($1, $2)",
			code, longURL,
		)
		if err == nil {
			return code, nil
		}

		// 23505 = unique_violation: the code is taken, so try another one
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			continue
		}
		return "", err
	}
	return "", errors.New("could not generate a unique code")
}

// GetLongURL returns the long URL stored under code, or ErrNotFound
func GetLongURL(pool *pgxpool.Pool, ctx context.Context, code string) (string, error) {
	var long_url string
	err := pool.QueryRow(ctx, `SELECT long_url FROM shortlinks WHERE code = $1`, code).Scan(&long_url)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", err
	}
	return long_url, nil
}

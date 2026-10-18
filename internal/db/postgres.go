package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"short.ly/internal/shortcode"
)

type DB struct {
	Pool *pgxpool.Pool
}

// Connect opens a connection pool to url and checks the database is reachable
func Connect(ctx context.Context, url string) (*DB, error) {
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
	return &DB{Pool: pool}, nil
}

func (d *DB) Close() {
	d.Pool.Close()
}

var ErrNotFound = errors.New("link not found")

// maxCodeAttempts bounds retries when a random code collides with an existing one.
const maxCodeAttempts = 5

// CreateLink saves a long URL under a freshly generated random code and returns the code
func (d *DB) CreateLink(ctx context.Context, longURL string) (string, error) {
	for range maxCodeAttempts {
		code, err := shortcode.GenerateCode(shortcode.CodeLength)
		if err != nil {
			return "", err
		}

		_, err = d.Pool.Exec(ctx,
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
func (d *DB) GetLongURL(ctx context.Context, code string) (string, error) {
	var long_url string
	err := d.Pool.QueryRow(ctx, `SELECT long_url FROM shortlinks WHERE code = $1`, code).Scan(&long_url)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", err
	}
	return long_url, nil
}

package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snake1-1eyes/Test_Ozon/internal/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ShortenerRepo struct {
	pool *pgxpool.Pool
}

func NewShortenerRepo(pool *pgxpool.Pool) *ShortenerRepo {
	return &ShortenerRepo{
		pool: pool,
	}
}

func NewPostgresPool(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %v", err)
	}

	cfg.MinConns = 5
	cfg.MaxConns = 30
	cfg.MaxConnLifetime = 3600
	cfg.MaxConnIdleTime = 1800
	cfg.HealthCheckPeriod = 300

	pool, err := pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %v", err)
	}

	return pool, nil
}

func (repo *ShortenerRepo) CreateAndSaveShortLink(originalURL string) (string, error) {
	var existingShortURL string
	err := repo.pool.QueryRow(
		context.Background(),
		"SELECT short_url FROM shorteners WHERE original_url = $1",
		originalURL,
	).Scan(&existingShortURL)
	if err == nil {
		return existingShortURL, nil
	}

	var shortURL string
	for i := 0; i < 5; i++ {
		shortURL, err = models.GenerateShortURL()
		if err != nil {
			return "", fmt.Errorf("failed to generate short URL: %v", err)
		}
		var dummy string
		err = repo.pool.QueryRow(
			context.Background(),
			"SELECT original_url FROM shorteners WHERE short_url = $1",
			shortURL,
		).Scan(&dummy)
		if err != nil {
			break
		}
		if i == 4 {
			return "", errors.New("failed to generate a unique short URL after 5 attempts")
		}
	}

	_, err = repo.pool.Exec(
		context.Background(),
		"INSERT INTO shorteners (original_url, short_url) VALUES ($1, $2)",
		originalURL, shortURL,
	)
	if err != nil {
		return "", fmt.Errorf("failed to save URL: %v", err)
	}

	return shortURL, nil
}

func (repo *ShortenerRepo) GetShortLink(shortURL string) (string, error) {
	var originalURL string
	err := repo.pool.QueryRow(
		context.Background(),
		"SELECT original_url FROM shorteners WHERE short_url = $1",
		shortURL,
	).Scan(&originalURL)
	if err != nil {
		return "", fmt.Errorf("failed to get original URL: %v", err)
	}
	return originalURL, nil
}

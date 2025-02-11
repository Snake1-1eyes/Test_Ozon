package postgres

import (
	"context"
	"errors"

	"github.com/Snake1-1eyes/Test_Ozon/internal/models"
	"github.com/jackc/pgtype/pgxtype"
	"github.com/jackc/pgx/v4"
)

type ShortenerRepo struct {
	db   pgxtype.Querier
	conn pgx.Conn
}

func NewShortenerRepo(db pgxtype.Querier, conn pgx.Conn) *ShortenerRepo {
	return &ShortenerRepo{
		db:   db,
		conn: conn,
	}
}

func (repo *ShortenerRepo) CreateAndSaveShortLink(originalURL string) (string, error) {
	var existingShortURL string
	err := repo.db.QueryRow(
		context.TODO(),
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
			return "", err
		}
		var dummy string
		err = repo.db.QueryRow(
			context.TODO(),
			"SELECT original_url FROM shorteners WHERE short_url = $1",
			shortURL,
		).Scan(&dummy)
		if err != nil {
			break
		}
		if i == 4 {
			return "", errors.New("failed to generate a unique short URL")
		}
	}

	_, err = repo.db.Exec(
		context.TODO(),
		"INSERT INTO shorteners (original_url, short_url) VALUES ($1, $2)",
		originalURL, shortURL,
	)
	if err != nil {
		return "", err
	}

	return shortURL, nil
}

func (repo *ShortenerRepo) GetShortLink(shortURL string) (string, error) {
	var originalURL string
	err := repo.db.QueryRow(
		context.TODO(),
		"SELECT original_url FROM shorteners WHERE short_url = $1",
		shortURL,
	).Scan(&originalURL)
	if err != nil {
		return "", err
	}
	return originalURL, nil
}

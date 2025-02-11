package inmemory

import (
	"errors"
	"sync"

	"github.com/Snake1-1eyes/Test_Ozon/internal/models"
)

type InMemoryRepo struct {
	originalToShort map[string]string
	shortToOriginal map[string]string
	mu              sync.Mutex
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		originalToShort: make(map[string]string),
		shortToOriginal: make(map[string]string),
	}
}

func (repo *InMemoryRepo) CreateAndSaveShortLink(originalURL string) (string, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if short, ok := repo.originalToShort[originalURL]; ok {
		return short, nil
	}

	var short string
	var err error
	for i := 0; i < 5; i++ {
		short, err = models.GenerateShortURL()
		if err != nil {
			return "", err
		}
		if _, exists := repo.shortToOriginal[short]; !exists {
			repo.originalToShort[originalURL] = short
			repo.shortToOriginal[short] = originalURL
			return short, nil
		}
	}
	return "", errors.New("failed to generate a unique short URL after 5 attempts")
}

func (repo *InMemoryRepo) GetShortLink(shortURL string) (string, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if original, exists := repo.shortToOriginal[shortURL]; exists {
		return original, nil
	}
	return "", errors.New("original URL not found")
}

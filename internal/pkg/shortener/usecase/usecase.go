package usecase

import (
	"github.com/Snake1-1eyes/Test_Ozon/internal/pkg/shortener"
)

type ShortenerUsecase struct {
	repo shortener.ShortenerRepo
}

func NewShortenerUsecase(repo shortener.ShortenerRepo) *ShortenerUsecase {
	return &ShortenerUsecase{
		repo: repo,
	}
}

func (u *ShortenerUsecase) CreateAndSaveShortLink(originalURL string) (string, error) {
	return u.repo.CreateAndSaveShortLink(originalURL)
}

func (u *ShortenerUsecase) GetShortLink(shortURL string) (string, error) {
	return u.repo.GetShortLink(shortURL)
}

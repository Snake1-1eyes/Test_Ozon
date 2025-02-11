package models

import (
	"crypto/rand"
	"errors"
	"math/big"
	"net/url"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
	"0123456789" +
	"_"

const shortURLLength = 10

type Shortener struct {
	OriginalURL string
	ShortURL    string
}

func GenerateShortURL() (string, error) {
	b := make([]byte, shortURLLength)
	for i := range b {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", errors.New("failed to generate random index: " + err.Error())
		}
		b[i] = charset[idx.Int64()]
	}
	return string(b), nil
}

func ValidateBaseURL(p *Shortener) error {
	parsedURL, err := url.ParseRequestURI(p.OriginalURL)
	if err != nil {
		return err
	}
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return errors.New("invalid URL: missing scheme or host")
	}
	return nil
}

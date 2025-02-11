package shortener

type ShortenerRepo interface {
	CreateAndSaveShortLink(originalURL string) (string, error)
	GetShortLink(shortURL string) (string, error)
}

type ShortenerUsecase interface {
	CreateAndSaveShortLink(originalURL string) (string, error)
	GetShortLink(shortURL string) (string, error)
}

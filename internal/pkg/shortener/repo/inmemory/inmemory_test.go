package inmemory

import (
	"testing"
)

func TestInMemoryRepo_CreateAndSaveShortLink(t *testing.T) {
	repo := NewInMemoryRepo()

	tests := []struct {
		name        string
		originalURL string
		wantErr     bool
	}{
		{
			name:        "Create new short URL",
			originalURL: "https://example.com",
			wantErr:     false,
		},
		{
			name:        "Get existing short URL",
			originalURL: "https://example.com",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.CreateAndSaveShortLink(tt.originalURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateAndSaveShortLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == "" {
				t.Error("CreateAndSaveShortLink() got empty short URL")
			}
		})
	}
}

func TestInMemoryRepo_GetShortLink(t *testing.T) {
	repo := NewInMemoryRepo()
	originalURL := "https://example.com"
	shortURL, _ := repo.CreateAndSaveShortLink(originalURL)

	tests := []struct {
		name     string
		shortURL string
		want     string
		wantErr  bool
	}{
		{
			name:     "Get existing URL",
			shortURL: shortURL,
			want:     originalURL,
			wantErr:  false,
		},
		{
			name:     "Get non-existing URL",
			shortURL: "nonexistent",
			want:     "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetShortLink(tt.shortURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetShortLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetShortLink() got = %v, want %v", got, tt.want)
			}
		})
	}
}

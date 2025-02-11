package usecase

import (
	"errors"
	"testing"
)

type mockRepo struct {
	createFunc func(string) (string, error)
	getFunc    func(string) (string, error)
}

func (m *mockRepo) CreateAndSaveShortLink(originalURL string) (string, error) {
	return m.createFunc(originalURL)
}

func (m *mockRepo) GetShortLink(shortURL string) (string, error) {
	return m.getFunc(shortURL)
}

func TestShortenerUsecase_CreateAndSaveShortLink(t *testing.T) {
	tests := []struct {
		name        string
		originalURL string
		mockResp    string
		mockErr     error
		want        string
		wantErr     bool
	}{
		{
			name:        "Success",
			originalURL: "https://example.com",
			mockResp:    "abc123",
			mockErr:     nil,
			want:        "abc123",
			wantErr:     false,
		},
		{
			name:        "Repository error",
			originalURL: "https://example.com",
			mockResp:    "",
			mockErr:     errors.New("repo error"),
			want:        "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockRepo{
				createFunc: func(s string) (string, error) {
					return tt.mockResp, tt.mockErr
				},
			}
			u := NewShortenerUsecase(mock)
			got, err := u.CreateAndSaveShortLink(tt.originalURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateAndSaveShortLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CreateAndSaveShortLink() got = %v, want %v", got, tt.want)
			}
		})
	}
}

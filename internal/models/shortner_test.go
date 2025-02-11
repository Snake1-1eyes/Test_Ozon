package models

import (
	"testing"
)

func TestGenerateShortURL(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Success generation",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateShortURL()
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateShortURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != shortURLLength {
				t.Errorf("GenerateShortURL() got length = %v, want %v", len(got), shortURLLength)
			}
		})
	}
}

func TestValidateBaseURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "Valid URL",
			url:     "https://example.com",
			wantErr: false,
		},
		{
			name:    "Invalid URL - no scheme",
			url:     "example.com",
			wantErr: true,
		},
		{
			name:    "Invalid URL - empty",
			url:     "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Shortener{
				OriginalURL: tt.url,
			}
			if err := ValidateBaseURL(s); (err != nil) != tt.wantErr {
				t.Errorf("ValidateBaseURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

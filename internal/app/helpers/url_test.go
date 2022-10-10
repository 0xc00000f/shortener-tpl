package helpers

import (
	"github.com/0xc00000f/shortener-tpl/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEncodeURL(t *testing.T) {
	tests := []struct {
		name    string
		letters int
	}{
		{
			name:    "6 letters url",
			letters: 6,
		},
		{
			name:    "72 letters url",
			letters: 72,
		},
		{
			name:    "0 letters url",
			letters: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := encodeURL(tt.letters)
			assert.Equal(t, len(url), tt.letters)
		})
	}
}

func TestEncodeURLWithDefaultSize(t *testing.T) {
	tests := []struct {
		name    string
		letters int
	}{
		{
			name:    "need letters url",
			letters: 6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := encodeURL(tt.letters)
			assert.Equal(t, len(url), tt.letters)
		})
	}
}

func TestEncodeAndStoreURL(t *testing.T) {

	var storage = storage.NewStorage()
	//storage.Set("ytAA2Z", "https://google.com")
	//storage.Set("hNaU8l", "https://dzen.ru/")

	tests := []struct {
		name    string
		baseURL string
	}{
		{
			name:    "positive #1",
			baseURL: "https://google.com",
		},
		{
			name:    "positive #2",
			baseURL: "https://dzen.ru/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encodedURL := EncodeAndStoreURL(tt.baseURL, storage)

			assert.Equal(t, len(encodedURL), 6)

			baseURLFromStorage, ok := storage.Get(encodedURL)
			require.True(t, ok)

			assert.Equal(t, baseURLFromStorage, tt.baseURL)
		})
	}
}

func TestDecodeURLFromStorage(t *testing.T) {

	var storage = storage.NewStorage()
	storage.Set("ytAA2Z", "https://google.com")
	storage.Set("hNaU8l", "https://dzen.ru/")

	tests := []struct {
		name       string
		encodedURL string
		baseURL    string
		ok         bool
	}{
		{
			name:       "positive #1",
			encodedURL: "ytAA2Z",
			baseURL:    "https://google.com",
			ok:         true,
		},
		{
			name:       "positive #2",
			encodedURL: "hNaU8l",
			baseURL:    "https://dzen.ru/",
			ok:         true,
		},
		{
			name:       "negative #1",
			encodedURL: "not exist",
			baseURL:    "",
			ok:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseURL, ok := DecodeURLFromStorage(tt.encodedURL, storage)
			require.Equal(t, tt.ok, ok)

			assert.Equal(t, tt.baseURL, baseURL)

		})
	}
}

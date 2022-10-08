package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDataStorage_Get(t *testing.T) {

	var storage = dataStorage{}
	//storage.Set("Jjqtdk", "https://vk.com")
	storage.Set("ytAA2Z", "https://google.com")
	//storage.Set("4qwpBs", "https://onlyfans.com")
	storage.Set("hNaU8l", "https://dzen.ru/")

	tests := []struct {
		name   string
		key    string
		value  string
		exists bool
	}{
		{
			name:   "positive #1",
			key:    "ytAA2Z",
			value:  "https://google.com",
			exists: true,
		},
		{
			name:   "positive #2",
			key:    "hNaU8l",
			value:  "https://dzen.ru/",
			exists: true,
		},
		{
			name:   "negative #1",
			key:    "4qwpBs",
			value:  "",
			exists: false,
		},
		{
			name:   "negative #2",
			key:    "",
			value:  "",
			exists: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, ok := storage.Get(tt.key)

			assert.Equal(t, ok, tt.exists)
			assert.Equal(t, value, tt.value)
		})
	}
}

func TestDataStorage_Set(t *testing.T) {

	var storage = dataStorage{}
	storage.Set("Jjqtdk", "https://vk.com")
	storage.Set("ytAA2Z", "https://google.com")
	storage.Set("4qwpBs", "https://onlyfans.com")
	storage.Set("hNaU8l", "https://dzen.ru/")

	tests := []struct {
		name  string
		key   string
		value string
	}{
		{
			name:  "simple set",
			key:   "Jjqtdk",
			value: "https://vk.com",
		},
		{
			name:  "simple set #2",
			key:   "ytAA2Z",
			value: "https://google.com",
		},
		{
			name:  "rewrite key #1",
			key:   "Jjqtdk",
			value: "https://onlyfans.com",
		},
		{
			name:  "rewrite key with empty value #1",
			key:   "hNaU8l",
			value: "",
		},
		{
			name:  "empty key #1",
			key:   "",
			value: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			storage.Set(tt.key, tt.value)
			value, _ := storage.Get(tt.key)

			assert.Equal(t, value, tt.value)
		})
	}
}

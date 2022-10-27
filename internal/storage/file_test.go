package storage

import (
	"errors"
	"github.com/stretchr/testify/require"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileStorage_Get(t *testing.T) {

	file, err := os.CreateTemp(os.TempDir(), "testfilestorage*")
	require.NoError(t, err)

	storage, err := NewFileStorage(file.Name())
	require.NoError(t, err)

	storage.Store("ytAA2Z", "https://google.com")
	storage.Store("hNaU8l", "https://dzen.ru/")

	tests := []struct {
		name  string
		key   string
		value string
		err   error
	}{
		{
			name:  "positive #1",
			key:   "ytAA2Z",
			value: "https://google.com",
			err:   nil,
		},
		{
			name:  "positive #2",
			key:   "hNaU8l",
			value: "https://dzen.ru/",
			err:   nil,
		},
		{
			name:  "negative #1",
			key:   "4qwpBs",
			value: "",
			err:   errors.New("key is not exist"),
		},
		{
			name:  "negative #2",
			key:   "",
			value: "",
			err:   errors.New("empty string as a key isn't allowed"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := storage.Get(tt.key)

			assert.Equal(t, err, tt.err)
			assert.Equal(t, value, tt.value)
		})
	}
}

func TestFileStorage_Set(t *testing.T) {

	file, err := os.CreateTemp(os.TempDir(), "testfilestorage*")
	require.NoError(t, err)

	storage, err := NewFileStorage(file.Name())
	require.NoError(t, err)

	tests := []struct {
		name     string
		key      string
		value    string
		errStore error
		errGet   error
	}{
		{
			name:     "simple set",
			key:      "Jjqtdk",
			value:    "https://vk.com",
			errStore: nil,
			errGet:   nil,
		},
		{
			name:     "simple set #2",
			key:      "ytAA2Z",
			value:    "https://google.com",
			errStore: nil,
			errGet:   nil,
		},
		{
			name:     "empty value #1",
			key:      "hNaU8l",
			value:    "",
			errStore: errors.New("empty string as a value isn't allowed"),
			errGet:   errors.New("key is not exist"),
		},
		{
			name:     "empty key #1",
			key:      "",
			value:    "",
			errStore: errors.New("empty string as a key isn't allowed"),
			errGet:   errors.New("empty string as a key isn't allowed"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := storage.Store(tt.key, tt.value)
			assert.Equal(t, tt.errStore, err)

			value, err := storage.Get(tt.key)

			assert.Equal(t, tt.value, value)
			assert.Equal(t, tt.errGet, err)
		})
	}
}

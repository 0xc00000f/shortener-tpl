package storage

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
)

func TestMemoryStorage_Get(t *testing.T) {

	var storage = NewMemoryStorage(zap.L())
	require.NoError(t, storage.Store(uuid.Nil, "ytAA2Z", "https://google.com"))
	require.NoError(t, storage.Store(uuid.Nil, "hNaU8l", "https://dzen.ru/"))

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
			err:   ErrNoKeyFound,
		},
		{
			name:  "negative #2",
			key:   "",
			value: "",
			err:   ErrEmptyKey,
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

func TestMemoryStorage_Set(t *testing.T) {

	var storage = NewMemoryStorage(zap.L())

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
			errStore: ErrEmptyValue,
			errGet:   ErrNoKeyFound,
		},
		{
			name:     "empty key #1",
			key:      "",
			value:    "",
			errStore: ErrEmptyKey,
			errGet:   ErrEmptyKey,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := storage.Store(uuid.Nil, tt.key, tt.value)
			assert.Equal(t, tt.errStore, err)

			value, err := storage.Get(tt.key)

			assert.Equal(t, tt.value, value)
			assert.Equal(t, tt.errGet, err)
		})
	}
}

package storage_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"

	"github.com/0xc00000f/shortener-tpl/internal/storage"
)

func TestMemoryStorage_Get(t *testing.T) {
	t.Parallel()

	var memoryStorage = storage.NewMemoryStorage(zap.L())

	ctx := context.Background()

	require.NoError(t, memoryStorage.Store(ctx, uuid.Nil, "ytAA2Z", "https://google.com"))
	require.NoError(t, memoryStorage.Store(ctx, uuid.Nil, "hNaU8l", "https://dzen.ru/"))

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
			err:   storage.ErrNoKeyFound,
		},
		{
			name:  "negative #2",
			key:   "",
			value: "",
			err:   storage.ErrEmptyKey,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			value, err := memoryStorage.Get(ctx, tt.key)

			assert.Equal(t, err, tt.err)
			assert.Equal(t, value, tt.value)
		})
	}
}

func TestMemoryStorage_Set(t *testing.T) {
	t.Parallel()

	var memoryStorage = storage.NewMemoryStorage(zap.L())

	ctx := context.Background()

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
			errStore: storage.ErrEmptyValue,
			errGet:   storage.ErrNoKeyFound,
		},
		{
			name:     "empty key #1",
			key:      "",
			value:    "",
			errStore: storage.ErrEmptyKey,
			errGet:   storage.ErrEmptyKey,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := memoryStorage.Store(ctx, uuid.Nil, tt.key, tt.value)
			assert.Equal(t, tt.errStore, err)

			value, err := memoryStorage.Get(ctx, tt.key)

			assert.Equal(t, tt.value, value)
			assert.Equal(t, tt.errGet, err)
		})
	}
}

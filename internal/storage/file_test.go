package storage_test

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.uber.org/zap"

	"github.com/0xc00000f/shortener-tpl/internal/storage"
)

func TestFileStorage_Get(t *testing.T) {
	t.Parallel()

	file, err := os.CreateTemp(os.TempDir(), "testfilestorage*")
	require.NoError(t, err)

	fileStorage, err := storage.NewFileStorage(file.Name(), zap.L())
	require.NoError(t, err)

	ctx := context.Background()

	require.NoError(t, fileStorage.Store(ctx, uuid.Nil, "ytAA2Z", "https://google.com"))
	require.NoError(t, fileStorage.Store(ctx, uuid.Nil, "hNaU8l", "https://dzen.ru/"))

	tests := []struct {
		userID uuid.UUID
		name   string
		key    string
		value  string
		err    error
	}{
		{
			userID: uuid.Nil,
			name:   "positive #1",
			key:    "ytAA2Z",
			value:  "https://google.com",
			err:    nil,
		},
		{
			userID: uuid.Nil,
			name:   "positive #2",
			key:    "hNaU8l",
			value:  "https://dzen.ru/",
			err:    nil,
		},
		{
			userID: uuid.Nil,
			name:   "negative #1",
			key:    "4qwpBs",
			value:  "",
			err:    storage.ErrNoKeyFound,
		},
		{
			userID: uuid.Nil,
			name:   "negative #2",
			key:    "",
			value:  "",
			err:    storage.ErrEmptyKey,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			value, err := fileStorage.Get(ctx, tt.key)

			assert.Equal(t, err, tt.err)
			assert.Equal(t, value, tt.value)
		})
	}
}

func TestFileStorage_Set(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	file, err := os.CreateTemp(os.TempDir(), "testfilestorage*")
	require.NoError(t, err)

	fileStorage, err := storage.NewFileStorage(file.Name(), zap.L())
	require.NoError(t, err)

	tests := []struct {
		userID   uuid.UUID
		name     string
		key      string
		value    string
		errStore error
		errGet   error
	}{
		{
			userID:   uuid.Nil,
			name:     "simple set",
			key:      "Jjqtdk",
			value:    "https://vk.com",
			errStore: nil,
			errGet:   nil,
		},
		{
			userID:   uuid.Nil,
			name:     "simple set #2",
			key:      "ytAA2Z",
			value:    "https://google.com",
			errStore: nil,
			errGet:   nil,
		},
		{
			userID:   uuid.Nil,
			name:     "empty value #1",
			key:      "hNaU8l",
			value:    "",
			errStore: storage.ErrEmptyValue,
			errGet:   storage.ErrNoKeyFound,
		},
		{
			userID:   uuid.Nil,
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

			err := fileStorage.Store(ctx, tt.userID, tt.key, tt.value)
			assert.Equal(t, tt.errStore, err)

			value, err := fileStorage.Get(ctx, tt.key)

			assert.Equal(t, tt.value, value)
			assert.Equal(t, tt.errGet, err)
		})
	}
}

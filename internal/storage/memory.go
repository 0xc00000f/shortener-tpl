package storage

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/0xc00000f/shortener-tpl/internal/log"
	"github.com/0xc00000f/shortener-tpl/internal/models"
)

var (
	ErrNoKeyFound = errors.New("key is not exist")
	ErrEmptyKey   = errors.New("empty string as a key isn't allowed")
	ErrEmptyValue = errors.New("empty string as a value isn't allowed")
)

type MemoryStorage struct {
	storage map[string]models.URL
	history map[uuid.UUID]map[string]models.URL

	mu sync.RWMutex
	l  *zap.Logger
}

func NewMemoryStorage(logger *zap.Logger) *MemoryStorage {
	return &MemoryStorage{
		storage: make(map[string]models.URL),
		history: make(map[uuid.UUID]map[string]models.URL),
		l:       logger,
	}
}

//revive:disable-next-line
func (ms *MemoryStorage) Get(ctx context.Context, short string) (long string, err error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	ms.l.Info("input", zap.String("short", short))

	if len(short) == 0 {
		err = ErrEmptyKey
		return "", err
	}

	url, ok := ms.storage[short]
	if !ok {
		return "", ErrNoKeyFound
	}

	long = url.Long

	if !url.IsActive {
		return url.Long, URLDeletedError{}
	}

	ms.l.Info("function result", zap.String("long", long), zap.Error(err))

	return long, nil
}

//revive:disable-next-line
func (ms *MemoryStorage) Store(ctx context.Context, userID uuid.UUID, short, long string) (err error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.l.Info("input",
		zap.String("userID", userID.String()),
		zap.String("short", short),
		zap.String("long", long),
		zap.Error(err),
	)

	if len(short) == 0 {
		return ErrEmptyKey
	}

	if len(long) == 0 {
		return ErrEmptyValue
	}

	if userID != uuid.Nil {
		if _, ok := ms.history[userID]; !ok {
			ms.history[userID] = map[string]models.URL{}
		}

		ms.history[userID][short] = models.URL{
			UserID:   userID,
			Short:    short,
			Long:     long,
			IsActive: true,
		}
	}

	ms.storage[short] = models.URL{
		UserID:   userID,
		Short:    short,
		Long:     long,
		IsActive: true,
	}

	return nil
}

//revive:disable-next-line
func (ms *MemoryStorage) IsKeyExist(ctx context.Context, short string) (bool, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	_, ok := ms.storage[short]

	return ok, nil
}

//revive:disable-next-line
func (ms *MemoryStorage) GetAll(ctx context.Context, userID uuid.UUID) (result map[string]string, err error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	ms.l.Info("function input", zap.String("userID", userID.String()))
	almostResult := ms.history[userID]

	result = make(map[string]string)

	for short, url := range almostResult {
		result[short] = url.Long
	}

	ms.l.Info("function result", log.MapToFields(result)...)

	return result, nil
}

//revive:disable-next-line
func (ms *MemoryStorage) GetStats(ctx context.Context) (models.Stats, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	count := 0
	for _, kvpairs := range ms.history {
		count += len(kvpairs)
	}

	return models.Stats{CountUsers: len(ms.history), CountURLs: count}, nil
}

//revive:disable-next-line
func (ms *MemoryStorage) Delete(ctx context.Context, data []models.URL) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	for _, url := range data {
		ms.storage[url.Short] = url
		ms.history[url.UserID][url.Short] = url
	}

	return nil
}

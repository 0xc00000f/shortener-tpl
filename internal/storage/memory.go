package storage

import (
	"errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrNoKeyFound = errors.New("key is not exist")
	ErrEmptyKey   = errors.New("empty string as a key isn't allowed")
	ErrEmptyValue = errors.New("empty string as a value isn't allowed")
)

type MemoryStorage struct {
	storage map[string]string
	history map[uuid.UUID]map[string]string
	l       *zap.Logger
}

func NewMemoryStorage(logger *zap.Logger) MemoryStorage {
	return MemoryStorage{
		storage: make(map[string]string),
		history: make(map[uuid.UUID]map[string]string),
		l:       logger,
	}
}

func (ms MemoryStorage) Get(short string) (long string, err error) {
	if len(short) == 0 {
		err = ErrEmptyKey
		return "", err
	}
	long, ok := ms.storage[short]
	if !ok {
		return "", ErrNoKeyFound
	}
	return long, nil
}

func (ms MemoryStorage) Store(id uuid.UUID, short, long string) error {
	if len(short) == 0 {
		return ErrEmptyKey
	}
	if len(long) == 0 {
		return ErrEmptyValue
	}
	if id != uuid.Nil {
		ms.history[id][short] = long
	}
	ms.storage[short] = long
	return nil
}

func (ms MemoryStorage) IsKeyExist(short string) (bool, error) {
	_, ok := ms.storage[short]
	return ok, nil
}

func (ms MemoryStorage) GetAll(uuid uuid.UUID) (result map[string]string, err error) {
	return ms.history[uuid], nil
}

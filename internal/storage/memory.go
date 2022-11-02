package storage

import (
	"errors"

	"go.uber.org/zap"
)

var (
	ErrNoKeyFound = errors.New("key is not exist")
	ErrEmptyKey   = errors.New("empty string as a key isn't allowed")
	ErrEmptyValue = errors.New("empty string as a value isn't allowed")
)

type MemoryStorage struct {
	storage map[string]string
	l       *zap.Logger
}

func NewMemoryStorage(logger *zap.Logger) MemoryStorage {
	return MemoryStorage{
		storage: make(map[string]string),
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

func (ms MemoryStorage) Store(short, long string) error {
	if len(short) == 0 {
		return ErrEmptyKey
	}
	if len(long) == 0 {
		return ErrEmptyValue
	}
	ms.storage[short] = long
	return nil
}

func (ms MemoryStorage) IsKeyExist(short string) (bool, error) {
	_, ok := ms.storage[short]
	return ok, nil
}

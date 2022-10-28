package storage

import (
	"errors"
	"go.uber.org/zap"
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
		err = errors.New("empty string as a key isn't allowed")
		return "", err
	}
	long, ok := ms.storage[short]
	if !ok {
		return "", errors.New("key is not exist")
	}
	return long, nil
}

func (ms MemoryStorage) Store(short, long string) error {
	if len(short) == 0 {
		return errors.New("empty string as a key isn't allowed")
	}
	if len(long) == 0 {
		return errors.New("empty string as a value isn't allowed")
	}
	ms.storage[short] = long
	return nil
}

func (ms MemoryStorage) IsKeyExist(short string) (bool, error) {
	_, ok := ms.storage[short]
	return ok, nil
}
